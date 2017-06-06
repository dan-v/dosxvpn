package dosxvpn

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

const (
	dropletImage  = "coreos-beta"
	dropletName   = "dosxvpn"
	dropletSize   = "512mb"
	dropletRegion = "sfo2"
)

type Droplet struct {
	DropletID   int
	IPv4Address string

	ssh *sshKeyPair
}

type Option func(*options)

func DropletName(name string) Option {
	return func(opt *options) {
		opt.dropletName = name
	}
}

func DropletRegion(region string) Option {
	return func(opt *options) {
		opt.dropletRegion = region
	}
}

type options struct {
	dropletName   string
	dropletRegion string
	dropletSize   string
}

func RemoveAllDroplets(token string) ([]string, error) {
	oauthClient := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
	client := godo.NewClient(oauthClient)

	droplets, _, err := client.Droplets.List(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	// attempt removal of all dosxvpn droplets
	removedDroplets := make([]string, 0)
	for _, droplet := range droplets {
		if strings.Contains(droplet.Name, "dosxvpn") {
			_, err := client.Droplets.Delete(context.TODO(), droplet.ID)
			if err != nil {
				log.Println("Failed to remove droplet", droplet.Name, err)
			}
			removedDroplets = append(removedDroplets, droplet.Name)
		}
	}
	sort.Strings(removedDroplets)

	// attempt removal of all dosxvpn firewalls
	firewalls, _, err := client.Firewalls.List(context.TODO(), nil)
	if err == nil {
		for _, firewall := range firewalls {
			if strings.Contains(firewall.Name, "dosxvpn") {
				_, err := client.Firewalls.Delete(context.TODO(), firewall.ID)
				if err != nil {
					log.Println("Failed to remove firewall", firewall.Name, err)
				}
			}
		}
	}

	return removedDroplets, nil
}

func Deploy(accessToken string, opts ...Option) (*Droplet, error) {
	opt := options{
		dropletName:   dropletName,
		dropletRegion: dropletRegion,
		dropletSize:   dropletSize,
	}
	for _, o := range opts {
		o(&opt)
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	))
	client := godo.NewClient(oauthClient)

	// Create keypair to connect to instance
	keypair, err := createSSHKeyPair()
	if err != nil {
		return nil, err
	}

	// Build user data to initialize the droplet
	userData, err := buildUserData(&opt, keypair)
	if err != nil {
		return nil, err
	}

	// Launch the DigitalOcean droplet.
	createRequest := &godo.DropletCreateRequest{
		Name:     opt.dropletName,
		Region:   opt.dropletRegion,
		Size:     opt.dropletSize,
		UserData: userData,
		Image: godo.DropletCreateImage{
			Slug: dropletImage,
		},
	}

	// Query all the SSH keys on the account so we can include them
	// in the droplet.
	ctx := context.TODO()
	sshKeys, _, err := client.Keys.List(ctx, nil)
	if err != nil {
		return nil, err
	}
	if len(sshKeys) == 0 {
		return nil, errors.New("Need at least one SSH key uploaded to your DigitalOcean account. Go add an SSH key: https://cloud.digitalocean.com/settings/security")
	}
	for _, key := range sshKeys {
		keyToAdd := godo.DropletCreateSSHKey{ID: key.ID}
		createRequest.SSHKeys = append(createRequest.SSHKeys, keyToAdd)
	}

	d, _, err := client.Droplets.Create(ctx, createRequest)
	if err != nil {
		return nil, err
	}

	droplet := &Droplet{
		DropletID: d.ID,
		ssh:       keypair,
	}

	// A just-created droplet won't have any of the network IP addresses
	// quite yet. We have to poll until the droplet is provisioned and
	// they're populated.
	for attempt := 1; droplet.IPv4Address == ""; attempt++ {
		time.Sleep(time.Duration(attempt) * time.Second) // linear backoff

		d, _, err := client.Droplets.Get(ctx, droplet.DropletID)
		if err != nil {
			return nil, err
		}

		for _, nv4 := range d.Networks.V4 {
			if nv4.IPAddress != "" {
				droplet.IPv4Address = nv4.IPAddress
			}
		}
		if attempt >= 10 {
			return nil, fmt.Errorf("timeout waiting for provisioning of droplet %d", droplet.DropletID)
		}
	}

	fwRequest := &godo.FirewallRequest{
		Name: opt.dropletName,
		InboundRules: []godo.InboundRule{
			{
				Protocol:  "tcp",
				PortRange: "22",
				Sources: &godo.Sources{
					Addresses: []string{"0.0.0.0/0", "::/0"},
				},
			},
			{
				Protocol:  "udp",
				PortRange: "500",
				Sources: &godo.Sources{
					Addresses: []string{"0.0.0.0/0", "::/0"},
				},
			},
			{
				Protocol:  "udp",
				PortRange: "4500",
				Sources: &godo.Sources{
					Addresses: []string{"0.0.0.0/0", "::/0"},
				},
			},
		},
		OutboundRules: []godo.OutboundRule{
			{
				Protocol: "icmp",
				Destinations: &godo.Destinations{
					Addresses: []string{"0.0.0.0/0", "::/0"},
				},
			},
			{
				Protocol:  "tcp",
				PortRange: "all",
				Destinations: &godo.Destinations{
					Addresses: []string{"0.0.0.0/0", "::/0"},
				},
			},
			{
				Protocol:  "udp",
				PortRange: "all",
				Destinations: &godo.Destinations{
					Addresses: []string{"0.0.0.0/0", "::/0"},
				},
			},
		},
		DropletIDs: []int{d.ID},
	}

	// Setup firewall
	_, _, err = client.Firewalls.Create(context.TODO(), fwRequest)
	if err != nil {
		return nil, err
	}

	return droplet, nil
}
