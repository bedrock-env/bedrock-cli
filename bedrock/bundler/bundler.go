package bundler

import (
	"fmt"
	"github.com/bedrock-env/bedrock-cli/bedrock/helpers"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

type Options struct {
	OverwriteFiles bool
	SkipFiles      bool
	BedrockDir     string
	BundlePath     string
}

func Bundle(options Options) bool {
	bundleConfig := strings.Builder{}
	var bundleSucceeded bool

	var extensions []Extension
	viper.UnmarshalKey("bundle", &extensions)

	fmt.Println(helpers.WarnStyleBold.MarginLeft(0).Render("=> Cleaning bundle..."))
	removeOldBundle(options)
	ensureBundleDir(options)

	fmt.Println(helpers.WarnStyleBold.MarginLeft(0).Render("=> Installing bundle..."))

	for _, extension := range extensions {
		fmt.Println(helpers.InfoStyleBold.Render(extension.Name))

		validationErrors := extension.Validate()
		if len(validationErrors) > 0 {
			fmt.Println(helpers.ErrorStyle.MarginLeft(0).Render("Invalid extension!"))
			for _, ve := range validationErrors {
				fmt.Println(helpers.ErrorStyle.MarginLeft(2).Render("- " + ve.Error()))
			}

			break
		}

		prepareErr := extension.Prepare(options)

		if prepareErr != nil {
			fmt.Println(helpers.ErrorStyle.MarginLeft(2).Render("Failed with:"))
			fmt.Println(helpers.ErrorStyle.MarginLeft(4).Render(prepareErr.Error()))

			break
		}

		setupSucceeded := extension.Setup(options)

		if !setupSucceeded {
			fmt.Println(helpers.ErrorStyle.MarginLeft(2).Render("Setup failed! Halting bundle install."))

			break
		}

		bundleConfig.WriteString(extension.SourcePath + "\n")
		bundleSucceeded = true
	}

	if !bundleSucceeded {
		return false
	}

	fmt.Println(helpers.WarnStyleBold.MarginLeft(0).Render("=> Writing bundle config..."))

	err := os.WriteFile(filepath.Join(options.BedrockDir, "bundle", "load"),
		[]byte(bundleConfig.String()), 0744)
	if err != nil {
		fmt.Println(err.Error())
	}

	return true
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
