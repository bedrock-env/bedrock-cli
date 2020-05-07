package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/bedrock-env/bedrock-cli/extensions"
	"github.com/bedrock-env/bedrock-cli/helpers"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
)

var OverwriteFiles bool

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the Bedrock extension bundle",
	Long: "Install the Bedrock extensions listed in the bundle file",
	Run: func(cmd *cobra.Command, args []string) {
		BundleInstall()
	},
}

func BundleInstall() {
	RemoveOldBundle()
	EnsureBundleDir()

	desiredBundle := LoadBundle()

	var bundle []extensions.Extension

	for _, e := range desiredBundle {
		result := e.Install(extensions.InstallOptions{OverwriteFiles: OverwriteFiles, BedrockDir: helpers.BedrockDir })

		if result {
			bundle = append(bundle, e)
		}
	}

	var bundleData string
	for _, e := range bundle {
		bundleData = fmt.Sprintf("%s%s\n", bundleData, e.BasePath)
	}

	ioutil.WriteFile(filepath.Join(helpers.BedrockDir, "bundle", "load"), []byte(bundleData), 0744)
}

func LoadBundle() []extensions.Extension {
	data, err := ioutil.ReadFile(filepath.Join(helpers.BedrockDir, "bundle.json"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var bundleConfig []map[string]string
	json.Unmarshal([]byte(data), &bundleConfig)

	fmt.Println(bundleConfig)
	var bundle []extensions.Extension
	for _, element := range bundleConfig {
		fmt.Println(element)
		extension := extensions.Extension{
			Name: element["name"],
			Branch: element["branch"],
			Git: element["git"],
			Tag: element["tag"],
			Ref: element["ref"],
			Path: element["path"],
		}.Load(packageManager, extensions.InstallOptions{BedrockDir: helpers.BedrockDir })
		bundle = append(bundle, extension)
	}

	fmt.Println(bundle)
	return bundle
}

func RemoveOldBundle() {
	bundlePath := filepath.Join(helpers.BedrockDir, "bundle")

	if !helpers.Exists(bundlePath) {
		return
	}

	err := os.RemoveAll(bundlePath)
	if err != nil {
		fmt.Println("Unable to remove old bundle at", bundlePath)
		fmt.Println(err)
		os.Exit(1)
	}
}
func EnsureBundleDir() {
	bundlePath := filepath.Join(helpers.BedrockDir, "bundle")

	if !helpers.Exists(bundlePath) {
		_ = os.Mkdir(bundlePath, 0744)
	}
}

func init() {
	bundleCmd.AddCommand(installCmd)
	installCmd.PersistentFlags().BoolVar(&OverwriteFiles, "overwrite",
		false, "Overwrite existing files")
}
