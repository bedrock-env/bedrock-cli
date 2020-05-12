package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bedrock-env/bedrock-cli/bedrock/bundler"

	"github.com/bedrock-env/bedrock-cli/bedrock/helpers"
)

var OverwriteFiles bool

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the Bedrock extension bundler",
	Long: "Install the Bedrock extensions listed in the bundler file",
	Run: func(cmd *cobra.Command, args []string) {
		bundler.Install(bundler.Options{
			BedrockDir: helpers.BedrockDir,
			PackageManager: PackageManager,
			OverwriteFiles: OverwriteFiles,
		})
	},
}

func init() {
	bundleCmd.AddCommand(installCmd)
	installCmd.PersistentFlags().BoolVar(&OverwriteFiles, "overwrite",
		false, "Overwrite existing files")
}
