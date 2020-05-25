package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/bedrock-env/bedrock-cli/bedrock"
)

var bundleCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Manage the extension bundler",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := bedrock.Preflight()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(bundleCmd)
}
