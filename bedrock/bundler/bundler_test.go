package bundler

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveOldBundle(t *testing.T) {
	parentDir := os.TempDir()
	tmpDir, err := ioutil.TempDir(parentDir, "bedrock-*")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(tmpDir)

	bundlePath := filepath.Join(tmpDir, "bundle")
	os.Mkdir(bundlePath, 0744)

	opts := Options{BedrockDir: tmpDir}
	removeOldBundle(opts)

	_, err = os.Stat(bundlePath)
	if !os.IsNotExist(err) {
		log.Println("Bundle directory was not removed")
		os.RemoveAll(tmpDir)
		os.Exit(1)
	}
}
