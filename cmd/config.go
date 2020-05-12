package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bedrock-env/bedrock-cli/bedrock/helpers"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or manage the Bedrock configuration",
	Run: func(cmd *cobra.Command, args []string) {
		viper.WriteConfigAs(filepath.Join(helpers.Home, configFileName))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
