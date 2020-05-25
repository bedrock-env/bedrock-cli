package cmd

import (
	"github.com/spf13/cobra"
)

var Interactive bool

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Manage Bedrock installation and updates",
}

func init() {
	rootCmd.AddCommand(setupCmd)
	setupUpdateCmd.PersistentFlags().BoolVarP(&Interactive, "interactive", "i", true, "update interactivity")
}
