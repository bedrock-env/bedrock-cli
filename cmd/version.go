package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VERSION = "0.0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show Bedrock version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), "Bedrock", VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
