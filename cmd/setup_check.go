package cmd

import (
	"fmt"
	"github.com/bedrock-env/bedrock-cli/bedrock"
	"github.com/bedrock-env/bedrock-cli/bedrock/helpers"
	"github.com/spf13/cobra"
)

var setupCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check Bedrock setup",
	Run: func(cmd *cobra.Command, args []string) {
		coreResult := bedrock.CheckCore()

		if !coreResult.Found {
			fmt.Printf("%s\u0078%s Bedrock Core was not detected\n", helpers.ColorRed, helpers.ColorReset)
		} else if !coreResult.MeetsMinVersion {
			fmt.Printf("%s\u0078%s Bedrock Core %s found, %s or newer required\n", helpers.ColorRed,
				coreResult.Version, bedrock.CoreMinVersion, helpers.ColorReset)
		} else if coreResult.UpdateAvailable {
			fmt.Printf("%s\u26A0%s An update is available for Bedrock Core \n", helpers.ColorYellow, helpers.ColorReset)
		} else {
			fmt.Printf("%s\u2714%s Bedrock Core detected\n", helpers.ColorGreen, helpers.ColorReset)
		}

		bedrock.CheckZSH()
		bedrock.CheckGit()
	},
}

func init() {
	setupCmd.AddCommand(setupCheckCmd)
}
