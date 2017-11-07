package web

import (
	"bytes"
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
	"path/filepath"
	"time"

	"github.com/dan-v/dosxvpn/deploy"
	"github.com/dan-v/dosxvpn/doclient"
)

type handler struct {
	oAuthClientID string
	token         string
	state         string
	host          string
	indexTmpl     *template.Template
	callbackTmpl  *template.Template
	dashboardTmpl *template.Template
	deleteTmpl    *template.Template
	regionTmpl    *template.Template
	progressTmpl  *template.Template
	completeTmpl  *template.Template
	deployment    *deploy.Deployment
}

func webHandler(host, oAuthClientID string) http.Handler {
	h := &handler{
		host:          host,
		oAuthClientID: oAuthClientID,
		indexTmpl:     template.Must(template.New("index page").Parse(indexPageHTML)),
		callbackTmpl:  template.Must(template.New("callback").Parse(callbackHTML)),
		dashboardTmpl: template.Must(template.New("dashboard").Parse(dashboardPageHTML)),
		deleteTmpl:    template.Must(template.New("delete").Parse(deletePageHTML)),
		progressTmpl:  template.Must(template.New("progress").Parse(progressPageHTML)),
		completeTmpl:  template.Must(template.New("complete").Parse(completePageHTML)),
	}
	h.generateState()
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", h.index)
	mux.HandleFunc("/callback", h.callback)
	mux.HandleFunc("/dashboard", h.dashboardPage)
	mux.HandleFunc("/delete", h.delete)
	mux.HandleFunc("/install", h.progressPage)
	mux.HandleFunc("/status/", h.status)
	mux.HandleFunc("/complete", h.completePage)
	mux.HandleFunc("/download", h.download)
	mux.HandleFunc("/exit", h.exit)
	return mux
}

func (h *handler) index(rw http.ResponseWriter, req *http.Request) {
	vals := make(url.Values)
	vals.Set("response_type", "token")
	vals.Set("client_id", h.oAuthClientID)
	vals.Set("state", h.state)
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

func (h *handler) dashboardPage(rw http.ResponseWriter, req *http.Request) {
	if h.token == "" {
		token, state := req.FormValue("access_token"), req.FormValue("state")
		if token == "" || state != h.state {
			http.Error(rw, "invalid oauth2 grant", http.StatusBadRequest)
			return
		}
		h.token = token
	}

	client := doclient.New(h.token)
	regions, err := client.ListRegions()
	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to list regions: %v", err), http.StatusBadRequest)
		return
	}

	vpnList, err := deploy.ListVpns(h.token)
	if err != nil {
		http.Error(rw, "failed to list droplets", http.StatusBadRequest)
		return
	}

	tmplData := struct {
		Regions map[string]string
		VPNList []string
	}{
		Regions: regions,
		VPNList: vpnList,
	}

	err = h.dashboardTmpl.Execute(rw, tmplData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "executing template: %s", err.Error())
	}
}

func (h *handler) delete(rw http.ResponseWriter, req *http.Request) {
	droplet := req.FormValue("droplet")
	if droplet == "" {
		rw.Write([]byte("Need to specify droplet"))
		return
	}
	_, err := deploy.RemoveVPN(h.token, droplet, true)
	if err != nil {
		rw.Write([]byte(fmt.Sprintf("Failed to remove VPN: %v", err)))
		return
	}

	remaining := false
	for i := 0; i < 5; i++ {
		vpns, _ := deploy.ListVpns(h.token)
		if err != nil {
			for _, vpn := range vpns {
				if vpn == droplet {
					remaining = true
				}
			}
		}
		if !remaining {
			break
		}
	}
	if !remaining {
		rw.Write([]byte("success"))
	} else {
		rw.Write([]byte("failure"))
	}
}

func (h *handler) region(rw http.ResponseWriter, req *http.Request) {
	token, state := req.FormValue("access_token"), req.FormValue("state")
	if token == "" || state != h.state {
		http.Error(rw, "invalid oauth2 grant", http.StatusBadRequest)
		return
	}
	h.token = token

	client := doclient.New(token)
	regions, err := client.ListRegions()
	if err != nil {
		http.Error(rw, "failed to list regions", http.StatusBadRequest)
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

func (h *handler) progressPage(rw http.ResponseWriter, req *http.Request) {
	region := req.FormValue("region")
	var err error
	h.deployment, err = deploy.New(h.token, region, true)
	if err != nil {
		return
	}
	go h.deployment.Run()

	err = h.progressTmpl.Execute(rw, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "executing template: %s", err.Error())
	}
}

func (h *handler) status(rw http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(h.deployment)

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(buf.Bytes())
}

func (h *handler) completePage(rw http.ResponseWriter, req *http.Request) {
	tmplData := struct {
		FinalIP  string
		Password string
	}{
		FinalIP:  h.deployment.VPNIPAddress,
		Password: h.deployment.VpnPassword,
	}

	err := h.completeTmpl.Execute(rw, tmplData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "executing template: %s", err.Error())
	}
}

func (h *handler) download(rw http.ResponseWriter, req *http.Request) {
	fileType := req.FormValue("type")
	name := h.deployment.Name
	if fileType == "apple" {
		savePath := filepath.Join(deploy.FilepathDosxvpnConfigDir, fmt.Sprintf(deploy.FilenameAppleConfig, name))
		data, err := ioutil.ReadFile(savePath)
		if err != nil {
			log.Println(err)
		}
		http.ServeContent(rw, req, name+".apple.mobileconfig", time.Now(), bytes.NewReader(data))
	} else if fileType == "android" {
		savePath := filepath.Join(deploy.FilepathDosxvpnConfigDir, fmt.Sprintf(deploy.FilenameAndroidConfig, name))
		data, err := ioutil.ReadFile(savePath)
		if err != nil {
			log.Println(err)
		}
		http.ServeContent(rw, req, name+".android.sswan", time.Now(), bytes.NewReader(data))
	} else {
		http.Error(rw, "invalid type to download", http.StatusBadRequest)
	}
}

func (h *handler) generateState() error {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}
	h.state = hex.EncodeToString(b)
	return nil
}

func (h *handler) exit(rw http.ResponseWriter, req *http.Request) {
	os.Exit(0)
}
