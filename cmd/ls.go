package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/dan-v/dosxvpn/deploy"
)

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "List dosxvpn VPN servers",
	Args: func(cmd *cobra.Command, args []string) error {
		if !digitalOceanTokenEnvSet() {
			return errorMissingToken
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		token := getCliToken()
		droplets, err := deploy.ListVpns(token)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(droplets)
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
