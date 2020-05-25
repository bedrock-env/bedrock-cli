package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bedrock-env/bedrock-cli/bedrock"
	"github.com/bedrock-env/bedrock-cli/bedrock/helpers"
)

var setupInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Bedrock Core",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: check if Bedrock Core is already installed
		ok, err := bedrock.InstallCore(Interactive)
		if ok {
			fmt.Printf("%s\u2714%s Bedrock Core was installed successfully\n", helpers.ColorGreen, helpers.ColorReset)
			fmt.Printf("    %s%s\n\n    %s%s\n",
				helpers.ColorYellow,
				`
    =========================
    Post-install instructions
    =========================
`,
				"Add the following to $HOME/.zshrc: source $HOME/.bedrock/bedrock.zsh",
				helpers.ColorReset)
		} else {
			fmt.Printf("%s\u0078%s Bedrock Core installation failed\n", helpers.ColorRed, helpers.ColorReset)
			fmt.Println(err)
		}
	},
}

func init() {
	setupCmd.AddCommand(setupInstallCmd)
}
