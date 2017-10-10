package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "dosxvpn",
	Short: "One click personal VPN server on DigitalOcean",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	fmt.Println(`
___  ____ ____ _  _ _  _ ___  _  _ 
|  \ |  | [__   \/  |  | |__] |\ | 
|__/ |__| ___] _/\_  \/  |    | \| 

`)
}
