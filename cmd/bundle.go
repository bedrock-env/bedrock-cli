package cmd

import (
	"github.com/spf13/cobra"
)

var bundleCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Manage the extension bundler",
}

func init() {
	rootCmd.AddCommand(bundleCmd)
}
