package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/dan-v/dosxvpn/deploy"
)

var name string
var removeProfile bool

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove dosxvpn VPN server",
	Args: func(cmd *cobra.Command, args []string) error {
		if name == "" {
			return errorMissingName
		}
		if !digitalOceanTokenEnvSet() {
			return errorMissingToken
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		_, err := deploy.RemoveVPN(getCliToken(), name, removeProfile)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Successfully removed", name)
	},
}

func init() {
	RootCmd.AddCommand(rmCmd)
	rmCmd.Flags().StringVar(&name, "name", "", "Name of droplet to remove")
	rmCmd.Flags().BoolVar(&removeProfile, "remove-profile", false, "Remove VPN profile as well (only for OSX).")
}
