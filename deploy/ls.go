package deploy

import (
	"sort"
	"strings"

	"github.com/dan-v/dosxvpn/doclient"
)

func ListVpns(token string) ([]string, error) {
	client := doclient.New(token)
	droplets, err := client.ListDroplets()
	if err != nil {
		return nil, err
	}

	allDroplets := make([]string, 0)
	for _, droplet := range droplets {
		if strings.Contains(droplet.Name, DropletBaseName) {
			allDroplets = append(allDroplets, droplet.Name)
		}
	}
	sort.Strings(allDroplets)
	return allDroplets, nil
}
