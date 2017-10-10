package cmd

import (
	"github.com/spf13/cobra"

	"github.com/dan-v/dosxvpn/web"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Run web installer",

	Run: func(cmd *cobra.Command, args []string) {
		cleanup := false
		if len(args) > 0 {
			cleanup = true
		}
		web.Run(cleanup)
	},
}

func init() {
	RootCmd.AddCommand(webCmd)
}
