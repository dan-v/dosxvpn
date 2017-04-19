package dosxvpn

import (
	"context"
	"fmt"
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
	return droplet, nil
}
