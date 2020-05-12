package bundler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bedrock-env/bedrock-cli/bedrock/helpers"
)

type Options struct {
	OverwriteFiles bool
	BedrockDir string
	PackageManager string
}

func Install(options Options) {
	removeOldBundle(options)
	ensureBundleDir(options)

	desiredBundle := load(options)

	var bundle []Extension

	for _, e := range desiredBundle {
		result := e.Install(options)

		if result {
			bundle = append(bundle, e)
		}
	}

	var bundleData string
	for _, e := range bundle {
		bundleData = fmt.Sprintf("%s%s\n", bundleData, e.BasePath)
	}

	ioutil.WriteFile(filepath.Join(options.BedrockDir, "bundle", "load"),
		[]byte(bundleData),0744)
}

func load(options Options) []Extension {
	data, err := ioutil.ReadFile(filepath.Join(options.BedrockDir, "bundle.json"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var bundleConfig []map[string]string
	json.Unmarshal([]byte(data), &bundleConfig)

	var bundle []Extension
	for _, element := range bundleConfig {
		extension := Extension{
			Name: element["name"],
			Branch: element["branch"],
			Git: element["git"],
			Tag: element["tag"],
			Ref: element["ref"],
			Path: element["path"],
		}.Init(options)
		bundle = append(bundle, extension)
	}

	return bundle
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
