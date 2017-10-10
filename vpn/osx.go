package vpn

import (
	"fmt"
	"os/exec"
	"os/user"
)

func OSXAddVPN(mobileConfigPath string) error {
	user, err := user.Current()
	if err != nil {
		return err
	}
	cmd := exec.Command("profiles", "-U", user.Username, "-I", "-F", mobileConfigPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to add VPN. Output: %s. Error: %v", output, err)
	}

	exec.Command("open", "/System/Library/CoreServices/Menu Extras/VPN.menu/").Start()

	return nil
}

func OSXRemoveVPN(name string) error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	cmd := exec.Command("profiles", "-U", user.Username, "-R", "-p", "com.github.dan-v.dosxvpn."+name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to remove VPN. Output: %s. Error: %v", output, err)
	}
	return nil
}
