package cmd

import (
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Bedrock Core",
	//Run: func(cmd *cobra.Command, args []string) {
	//
	//}
}

func init() {
	rootCmd.AddCommand(bundleCmd)
}
