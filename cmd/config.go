package cmd

import (
	"github.com/bedrock-env/bedrock-cli/bedrock/helpers"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage the Bedrock configuration",
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringVarP(&helpers.BedrockDir, "bedrock-dir", "", filepath.Join(helpers.Home, ".bedrock"), "Set the Bedrock base directory. (absolute path)")
	viper.BindPFlag("settings.bedrock_dir", configCmd.LocalFlags().Lookup("bedrock-dir"))
}
