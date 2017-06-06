package dosxvpn

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

func Handler(oauthClientID, host string) http.Handler {
	h := &handler{
		oauthClientID: oauthClientID,
		host:          host,
		progressTmpl:  template.Must(template.New("progress").Parse(progressPageHTML)),
		callbackTmpl:  template.Must(template.New("callback").Parse(callbackHTML)),
		indexTmpl:     template.Must(template.New("index").Parse(indexPageHTML)),
		regionTmpl:    template.Must(template.New("region").Parse(regionPageHTML)),
		uninstallTmpl: template.Must(template.New("region").Parse(uninstallPageHTML)),
		installs:      make(map[string]*install),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.index)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/callback", h.callback)
	mux.HandleFunc("/grant", h.grant)
	mux.HandleFunc("/region", h.region)
	mux.HandleFunc("/install/", h.progressPage)
	mux.HandleFunc("/status/", h.status)
	mux.HandleFunc("/uninstall", h.uninstall)
	mux.HandleFunc("/exit", h.exit)
	mux.HandleFunc("/download", h.download)
	return mux
}

type handler struct {
	oauthClientID string
	host          string
	progressTmpl  *template.Template
	callbackTmpl  *template.Template
	indexTmpl     *template.Template
	regionTmpl    *template.Template
	uninstallTmpl *template.Template

	installMu sync.Mutex
	installs  map[string]*install
}

func (h *handler) index(rw http.ResponseWriter, req *http.Request) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	installID := hex.EncodeToString(b)

	h.installMu.Lock()
	h.installs[installID] = &install{Status: "pending auth"}
	h.installMu.Unlock()

	vals := make(url.Values)
	vals.Set("response_type", "token")
	vals.Set("client_id", h.oauthClientID)
	vals.Set("state", installID)
	vals.Set("scope", "read write")
	vals.Set("redirect_uri", h.host+"/callback")
	u := url.URL{
		Scheme:   "https",
		Host:     "cloud.digitalocean.com",
		Path:     "/v1/oauth/authorize",
		RawQuery: vals.Encode(),
	}

	tmplData := struct {
		InstallLink string
	}{
		InstallLink: u.String(),
	}
	h.indexTmpl.Execute(rw, tmplData)
}

func (h *handler) callback(rw http.ResponseWriter, req *http.Request) {
	err := h.callbackTmpl.Execute(rw, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "executing template: %s", err.Error())
	}
}

func (h *handler) grant(rw http.ResponseWriter, req *http.Request) {
	token, state := req.FormValue("access_token"), req.FormValue("state")
	if token == "" || state == "" {
		http.Error(rw, "invalid oauth2 grant", http.StatusBadRequest)
		return
	}

	h.installMu.Lock()
	curr := h.installs[state]
	h.installMu.Unlock()
	if curr == nil {
		http.Error(rw, "invalid oauth2 state", http.StatusBadRequest)
		return
	}

	curr.mu.Lock()
	curr.accessToken = token
	curr.mu.Unlock()

	http.Redirect(rw, req, "/region?access_token="+token+"&state="+state, http.StatusFound)
}

func (h *handler) region(rw http.ResponseWriter, req *http.Request) {
	token, state := req.FormValue("access_token"), req.FormValue("state")
	if token == "" || state == "" {
		http.Error(rw, "invalid oauth2 grant", http.StatusBadRequest)
		return
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
	client := godo.NewClient(oauthClient)
	r, _, err := client.Regions.List(context.TODO(), nil)
	if err != nil {
		http.Error(rw, "Failed to get list of regions", http.StatusBadRequest)
		return
	}
	regions := make(map[string]string)
	for _, region := range r {
		regions[region.Slug] = region.Name
	}

	tmplData := struct {
		Regions map[string]string
	}{
		Regions: regions,
	}

	err = h.regionTmpl.Execute(rw, tmplData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "executing template: %s", err.Error())
	}
}

func (h *handler) uninstall(rw http.ResponseWriter, req *http.Request) {
	token := req.FormValue("access_token")
	if token == "" {
		http.Error(rw, "invalid oauth2 grant", http.StatusBadRequest)
		return
	}

	removedDroplets, err := RemoveAllDroplets(token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "executing template: %s", err.Error())
	}

	tmplData := struct {
		RemovedDroplets []string
	}{
		RemovedDroplets: removedDroplets,
	}

	err = h.uninstallTmpl.Execute(rw, tmplData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "executing template: %s", err.Error())
	}
}

func (h *handler) progressPage(rw http.ResponseWriter, req *http.Request) {
	region := req.FormValue("region")
	id := path.Base(req.URL.Path)
	h.installMu.Lock()
	curr := h.installs[id]
	h.installMu.Unlock()

	if curr == nil {
		http.NotFound(rw, req)
		return
	}

	go curr.init(id, region)

	tmplData := struct {
		InstallID string
	}{
		InstallID: id,
	}
	err := h.progressTmpl.Execute(rw, tmplData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "executing template: %s", err.Error())
	}
}

func (h *handler) status(rw http.ResponseWriter, req *http.Request) {
	id := path.Base(req.URL.Path)
	h.installMu.Lock()
	curr := h.installs[id]
	h.installMu.Unlock()

	if curr == nil {
		http.NotFound(rw, req)
		return
	}

	var buf bytes.Buffer
	curr.mu.Lock()
	_ = json.NewEncoder(&buf).Encode(curr)
	curr.mu.Unlock()

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(buf.Bytes())
}

func (h *handler) exit(rw http.ResponseWriter, req *http.Request) {
	os.Exit(0)
}

func (h *handler) download(rw http.ResponseWriter, req *http.Request) {
	data, err := ioutil.ReadFile("/tmp/dosxvpn.mobileconfig")
	if err != nil {
		log.Println(err)
	}
	http.ServeContent(rw, req, "dosxvpn.mobileconfig", time.Now(), bytes.NewReader(data))
}

type install struct {
	mu           sync.Mutex
	Status       string `json:"status"`
	VPNIPAddress string `json:"ip_address"`
	InitialIP    string `json:"initial_ip"`
	FinalIP      string `json:"final_ip"`
	accessToken  string
	c            *Droplet
}

func (i *install) setStatus(status string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.Status = status
}

func getPublicIp() (string, error) {
	log.Println("Getting public IP address..")
	resp, err := http.Get("http://checkip.amazonaws.com/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(buf)), nil
}

func (i *install) init(state, region string) {
	defer revoke(i.accessToken)
	var droplet *Droplet
	var err error

	defer func() {
		if err != nil {
			i.setStatus(err.Error())
		}
	}()

	initialIP, _ := getPublicIp()
	i.mu.Lock()
	i.InitialIP = initialIP
	i.mu.Unlock()

	// Start deploying and create the droplet.
	droplet, err = Deploy(i.accessToken, DropletName("dosxvpn-"+state[:6]+"-"+region), DropletRegion(region))
	if err != nil {
		return
	}

	i.mu.Lock()
	i.VPNIPAddress = droplet.IPv4Address
	i.c = droplet
	i.Status = "waiting for ssh"
	i.mu.Unlock()

	err = WaitForSSH(droplet)
	if err != nil {
		return
	}

	i.setStatus("configuring vpn")
	vpnDetails, err := GetVPNDetails(droplet, region)
	if err != nil {
		return
	}

	i.setStatus("adding vpn to osx")
	err = SetupVPN(vpnDetails)
	if err != nil {
		return
	}

	i.setStatus("waiting for ip address change")
	for j := 0; j < 10; j++ {
		time.Sleep(time.Second * 5)
		newIp, err := getPublicIp()
		if err == nil && newIp != "" && newIp != initialIP {
			i.FinalIP = newIp
			break
		}
	}

	i.mu.Lock()
	i.Status = "done"
	i.c = nil // garbage collect the SSH keys
	i.mu.Unlock()
}

func revoke(accessToken string) error {
	body := strings.NewReader(url.Values{"token": {accessToken}}.Encode())
	req, err := http.NewRequest("POST", "https://cloud.digitalocean.com/v1/oauth/revoke", body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("revoke endpoint returned %d status code", resp.StatusCode)
	}
	return nil
}
