package bundler

import (
	"fmt"
	"github.com/bedrock-env/bedrock-cli/bedrock/helpers"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type Options struct {
	OverwriteFiles bool
	SkipFiles      bool
	BedrockDir     string
	PackageManager string
}

func Bundle(options Options) {
	var extensions []Extension
	var extensionManifest map[string]interface{}
	viper.UnmarshalKey("extensions", &extensionManifest)

	for k, v := range extensionManifest {
		var extension Extension
		err := mapstructure.Decode(v, &extension)
		if err != nil {
			fmt.Println(err)
		}
		extension.Name = k
		extensions = append(extensions, extension)
	}

	fmt.Println(helpers.WarnStyleBold.MarginLeft(0).Render("=> Cleaning bundle..."))
	// removeOldBundle(options)
	// ensureBundleDir(options)
	fmt.Println(helpers.WarnStyleBold.MarginLeft(0).Render("=> Installing bundle..."))

	//bundleConfig := strings.Builder{}

	for _, extension := range extensions {
		fmt.Println(helpers.InfoStyleBold.Render(extension.Name))
		//fmt.Println(helpers.BasicStyle.MarginLeft(2).Render(extension.Url))

		//for _, e := range extensions {
		//	succeeded := e.Install(options)
		//
		//	if succeeded {
		//		bundleConfig.WriteString(e.BasePath + "\n")
		//	}
		//}
		//
		//os.WriteFile(filepath.Join(options.BedrockDir, "bundle", "load"),
		//	[]byte(bundleConfig.String()), 0744)
	}
}

func removeOldBundle(options Options) {
	bundlePath := filepath.Join(options.BedrockDir, "bundle")

	if !helpers.Exists(bundlePath) {
		return
	}

	err := os.RemoveAll(bundlePath)
	if err != nil {
		fmt.Println("Unable to remove old bundler at", bundlePath)
		fmt.Println(err)
		os.Exit(1)
	}
}
func ensureBundleDir(options Options) {
	bundlePath := filepath.Join(options.BedrockDir, "bundle")

	if !helpers.Exists(bundlePath) {
		_ = os.Mkdir(bundlePath, 0744)
	}
}
