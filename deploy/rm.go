package deploy

import (
	"log"
	"strings"

	"github.com/dan-v/dosxvpn/doclient"
	"github.com/dan-v/dosxvpn/vpn"
)

func RemoveVPN(token, name string) ([]string, error) {
	log.Printf("Listing droplets..")
	client := doclient.New(token)
	droplets, err := client.ListDroplets()
	if err != nil {
		return nil, err
	}

	log.Printf("Listing firewalls..")
	firewalls, err := client.ListFirewalls()
	if err == nil {
		for _, firewall := range firewalls {
			if strings.Contains(firewall.Name, name) {
				err := client.DeleteFirewall(firewall.ID)
				if err != nil {
					log.Println("Failed to remove firewall", firewall.Name, err)
				}
			}
		}
	}

	log.Printf("Removing droplet..")
	removedDroplets := make([]string, 0)
	for _, droplet := range droplets {
		if strings.Contains(droplet.Name, name) {
			err := client.DeleteDroplet(droplet.ID)
			if err != nil {
				log.Println("Failed to remove droplet", droplet.Name, err)
			}
			removedDroplets = append(removedDroplets, droplet.Name)
		}
	}

	log.Printf("Removing OSX VPN profile for %s", name)
	err = vpn.OSXRemoveVPN(name)
	if err != nil {
		log.Printf("Failed to remove OSX VPN profile for %s. %v", name, err)
	}
	log.Printf("Finished removing OSX VPN profile for %s", name)

	return removedDroplets, nil
}
