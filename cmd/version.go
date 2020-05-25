package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/bedrock-env/bedrock-cli/bedrock"
)

const VERSION = "0.0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show Bedrock version information",
	PreRun: func(cmd *cobra.Command, args []string) {
		err := bedrock.Preflight()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), "Bedrock", VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
