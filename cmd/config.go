package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bedrock-env/bedrock-cli/bedrock/helpers"
)

var packageManager string

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage the Bedrock configuration",
	Run: func(cmd *cobra.Command, args []string) {
		// viper.WriteConfigAs(filepath.Join(helpers.Home, ".bedrock.json"))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringVarP(&packageManager, "package-manager", "", helpers.DefaultPkgManager(), "Set desired package manager")
	configCmd.Flags().StringVarP(&helpers.BedrockDir, "bedrock-dir", "", filepath.Join(helpers.Home, ".bedrock"), "Set the Bedrock base directory. (absolute path)")
	viper.BindPFlag("settings.package_manager", configCmd.LocalFlags().Lookup("package-manager"))
	viper.BindPFlag("settings.bedrock_dir", configCmd.LocalFlags().Lookup("bedrock-dir"))
}
