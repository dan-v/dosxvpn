package cmd

import (
	"log"

	"github.com/dan-v/dosxvpn/deploy"
	"github.com/spf13/cobra"
)

var region string
var autoConfigure bool

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy dosxvpn VPN server",
	Args: func(cmd *cobra.Command, args []string) error {
		if region == "" {
			return errorMissingRegion
		}
		if !digitalOceanTokenEnvSet() {
			return errorMissingToken
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		deployment, err := deploy.New(getCliToken(), region, autoConfigure)
		if err != nil {
			log.Fatal("Deployment failed:", err)
		}
		err = deployment.Run()
		if err != nil {
			log.Fatal("Deployment failed:", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVar(&region, "region", "", "Region to deploy droplet (e.g. ams2,ams3,nyc1,nyc2,nyc3,sfo1,sfo2).")
	deployCmd.Flags().BoolVar(&autoConfigure, "auto-configure", false, "Auto configure VPN (only for OSX).")
}
