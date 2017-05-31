package dosxvpn

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func GetVPNDetails(c *Droplet, region string) (string, error) {
	const getVpnDetails = `
	until docker logs dosxvpn &>/dev/null; do sleep 2; done; sleep 5; docker exec dosxvpn cat /etc/ipsec.d/client.mobileconfig
	`

	session, err := connect(c.IPv4Address, c.ssh)
	if err != nil {
		return "", err
	}

	rOut, err := session.StdoutPipe()
	if err != nil {
		return "", err
	}
	rErr, err := session.StderrPipe()
	if err != nil {
		return "", err
	}
	combined := io.MultiReader(rOut, rErr)

	err = session.Start(getVpnDetails)

	var lines []string
	scanner := bufio.NewScanner(combined)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	session.Close()

	output := strings.Join(lines, "\n")
	output = strings.TrimSpace(output)
	vpnNameXml := "<string>dosxvpn-" + region + " (" + c.IPv4Address + ")</string>"
	output = strings.Replace(output, "<string>dosxvpn</string>", vpnNameXml, -1)

	return output, nil
}

func SetupVPN(vpnDetails string) error {
	err := ioutil.WriteFile("/tmp/dosxvpn.mobileconfig", []byte(vpnDetails), 0644)
	if err != nil {
		return err
	}

	cmd := "profiles"
	args := []string{"-I", "-F", "/tmp/dosxvpn.mobileconfig"}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	exec.Command("open", "/System/Library/CoreServices/Menu Extras/VPN.menu/").Start()

	return nil
}
