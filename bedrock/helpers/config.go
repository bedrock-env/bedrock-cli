package helpers

import (
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"

	yamlv3 "gopkg.in/yaml.v3"
)

// RewriteConfig rewrites the config with desired indentation.
func RewriteConfig() {
	configPath := filepath.Join(Home, ".config", "bedrock", "config.yaml")
	yfile, err := ioutil.ReadFile(configPath)

	if err != nil {
		log.Fatal(err)
	}

	configData := make(map[interface{}]interface{})

	err2 := yamlv3.Unmarshal(yfile, &configData)

	if err2 != nil {
		log.Fatal(err2)
	}

	var b bytes.Buffer
	yamlEncoder := yamlv3.NewEncoder(&b)
	yamlEncoder.SetIndent(2) // this is what you're looking for
	yamlEncoder.Encode(&configData)

	ioutil.WriteFile(configPath, b.Bytes(), 0664)
}
