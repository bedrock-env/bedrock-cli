package cmd

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bedrock-env/bedrock-cli/bedrock"
)

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
		doc := strings.Builder{}

		versionBlock := lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Render("CLI"),
				lipgloss.NewStyle().MarginLeft(5).Render(bedrock.VERSION),
			),
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Render("Core"),
				lipgloss.NewStyle().MarginLeft(4).Render(bedrock.CoreVersion().String()),
			),
		)

		doc.WriteString(versionBlock)

		fmt.Fprintln(cmd.OutOrStdout(), doc.String())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
