package doclient

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

type Client struct {
	token    string
	doClient *godo.Client
}

func New(token string) *Client {
	ctx := context.Background()
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).DialContext,
		},
	}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	client := &Client{
		token:    token,
		doClient: godo.NewClient(oauthClient),
	}
	return client
}

func (c *Client) WaitForDropletIP(dropletID int) (ip string, err error) {
	for attempt := 1; ip == ""; attempt++ {
		time.Sleep(time.Duration(attempt) * time.Second) // linear backoff

		droplet, _, err := c.doClient.Droplets.Get(context.TODO(), dropletID)
		if err != nil {
			return "", err
		}

		for _, nv4 := range droplet.Networks.V4 {
			if nv4.IPAddress != "" {
				ip = nv4.IPAddress
			}
		}
		if attempt >= 10 {
			return "", fmt.Errorf("Timeout waiting for provisioning of droplet %v", dropletID)
		}
	}
	return ip, nil
}

func (c *Client) CreateSSHKey(name, publicKey string) (id int, err error) {
	createRequest := &godo.KeyCreateRequest{
		Name:      "dosxvpn",
		PublicKey: publicKey,
	}
	key, _, err := c.doClient.Keys.Create(context.TODO(), createRequest)
	if err != nil {
		return 0, err
	}
	return key.ID, nil
}

func (c *Client) CreateDroplet(name, region, size, userData, image string) (id int, err error) {
	createRequest := &godo.DropletCreateRequest{
		Name:     name,
		Region:   region,
		Size:     size,
		UserData: userData,
		Image: godo.DropletCreateImage{
			Slug: image,
		},
		IPv6: true,
	}

	accountSSHKeys, err := c.GetAccountSSHKeys()
	if err != nil {
		return 0, err
	}
	for _, key := range accountSSHKeys {
		keyToAdd := godo.DropletCreateSSHKey{ID: key.ID}
		createRequest.SSHKeys = append(createRequest.SSHKeys, keyToAdd)
	}

	droplet, _, err := c.doClient.Droplets.Create(context.TODO(), createRequest)
	if err != nil {
		return 0, err
	}
	return droplet.ID, nil
}

func (c *Client) CreateFirewall(name string, dropletID int) error {
	fwRequest := &godo.FirewallRequest{
		Name:          name,
		DropletIDs:    []int{dropletID},
		InboundRules:  c.generateInboundFirewallRules(),
		OutboundRules: c.generateOutboundFirewallRules(),
	}
	_, _, err := c.doClient.Firewalls.Create(context.TODO(), fwRequest)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ListDroplets() (droplets []godo.Droplet, err error) {
	droplets, _, err = c.doClient.Droplets.List(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return droplets, nil
}

func (c *Client) ListRegions() (regions map[string]string, err error) {
	r, _, err := c.doClient.Regions.List(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	regions = make(map[string]string)
	for _, region := range r {
		regions[region.Slug] = region.Name
	}
	return regions, nil
}

func (c *Client) DeleteDroplet(dropletID int) error {
	_, err := c.doClient.Droplets.Delete(context.TODO(), dropletID)
	if err != nil {
		return fmt.Errorf("Failed to remove droplet %v. %v", dropletID, err)
	}
	return nil
}

func (c *Client) ListFirewalls() (firewalls []godo.Firewall, err error) {
	firewalls, _, err = c.doClient.Firewalls.List(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return firewalls, nil
}

func (c *Client) DeleteFirewall(firewallID string) error {
	_, err := c.doClient.Firewalls.Delete(context.TODO(), firewallID)
	if err != nil {
		return fmt.Errorf("Failed to remove droplet %v. %v", firewallID, err)
	}
	return nil
}

func (c *Client) GetAccountSSHKeys() ([]godo.Key, error) {
	keys, _, err := c.doClient.Keys.List(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (c *Client) generateInboundFirewallRules() []godo.InboundRule {
	return []godo.InboundRule{
		{
			Protocol: "icmp",
			Sources: &godo.Sources{
				Addresses: []string{"0.0.0.0/0", "::/0"},
			},
		},
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
	}
}
func (c *Client) generateOutboundFirewallRules() []godo.OutboundRule {
	return []godo.OutboundRule{
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
	}
}
