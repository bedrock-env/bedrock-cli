package cmd

import (
	"fmt"
	"github.com/bedrock-env/bedrock-cli/extensions"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var OverwriteFiles bool

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the Bedrock extension bundle",
	Long: "Install the Bedrock extensions listed in the bundle file",
	Run: func(cmd *cobra.Command, args []string) {
		bundle := LoadBundle()

		for _, e := range bundle {
			e.Install(extensions.InstallOptions{OverwriteFiles: OverwriteFiles, BedrockDir: BedrockDir })
		}
	},
}

func LoadBundle() []extensions.Extension {
	b, err := ioutil.ReadFile(filepath.Join(BedrockDir, "bundle"))
	if err != nil {
		fmt.Println(err)
	}

	bundleList := strings.Split(strings.TrimSpace(string(b)), "\n")
	var bundle []extensions.Extension

	for _, e := range bundleList {
		extension := extensions.Extension{ Name: e }.Load(packageManager)
		bundle = append(bundle, extension)
	}

	return bundle
}

func init() {
	bundleCmd.AddCommand(installCmd)
	installCmd.PersistentFlags().BoolVar(&OverwriteFiles, "overwrite",
		false, "Overwrite existing files")
}
