package cmd

import (
	"fmt"
	"github.com/bedrock-env/bedrock-cli/bedrock/bundler"
	"github.com/bedrock-env/bedrock-cli/bedrock/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var OverwriteFiles bool
var SkipFiles bool

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the Bedrock extension bundle",
	Long:  "Install the Bedrock extensions defined in the bundle file",
	Run: func(cmd *cobra.Command, args []string) {
		packageManager := viper.GetString("settings.package_manager")
		if len(packageManager) == 0 {
			packageManager = helpers.DefaultPkgManager()
		}

		if OverwriteFiles && SkipFiles {
			fmt.Println("WARN: overwrite-files and skip-files were set. Falling back to skip-files.")
			OverwriteFiles = false
		}

		bundlerOptions := bundler.Options{
			BedrockDir:     helpers.BedrockDir,
			PackageManager: packageManager,
			OverwriteFiles: OverwriteFiles,
			SkipFiles:      SkipFiles,
		}

		bundler.Bundle(bundlerOptions)
	},
}

func init() {
	bundleCmd.AddCommand(installCmd)
	installCmd.PersistentFlags().BoolVar(&OverwriteFiles, "overwrite-files",
		false, "Overwrite existing files during syncing")
	installCmd.PersistentFlags().BoolVar(&SkipFiles, "skip-files",
		false, "Skip syncing existing files")
}
