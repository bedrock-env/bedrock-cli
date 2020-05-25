package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bedrock-env/bedrock-cli/bedrock"
	"github.com/bedrock-env/bedrock-cli/bedrock/helpers"
)

var setupUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Bedrock Core",
	Run: func(cmd *cobra.Command, args []string) {
		ok, err := bedrock.UpdateCore(Interactive)
		if ok {
			fmt.Printf("%s\u2714%s Bedrock Core was updated successfully\n", helpers.ColorGreen, helpers.ColorReset)
		} else {
			fmt.Printf("%s\u0078%s Bedrock Core update failed\n", helpers.ColorRed, helpers.ColorReset)
			fmt.Println(err)
		}
	},
}

func init() {
	setupCmd.AddCommand(setupUpdateCmd)
}
