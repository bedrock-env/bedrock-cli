package extensions

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/bedrock-env/bedrock-cli/helpers"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Extension struct {
	Name string
	InstallSteps []InstallStep
	Files []File
	PostInstallMessages []string
	BasePath string
}

type Manifest struct {
	Author    string    `json:"author"`
	Platforms Platforms `json:"platforms"`
}
type InstallStep struct {
	Binary  string `json:"binary"`
	Command string `json:"command"`
	RunIf string   `json:"runif"`
}
type UpdateSteps struct {
	Binary  string `json:"binary"`
	Command string `json:"command"`
}
type Homebrew struct {
	InstallSteps []InstallStep `json:"install_steps"`
	UpdateSteps  []UpdateSteps  `json:"update_steps"`
	PostInstallMessages []string `json:"post_install_messages"`
}
type File struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Operation string `json:"operation"`
}
type Macos struct {
	Homebrew Homebrew `json:"homebrew"`
	Files    []File  `json:"files"`
}
type Apt struct {
	InstallSteps []InstallStep `json:"install_steps"`
	UpdateSteps  []UpdateSteps  `json:"update_steps"`
	PostInstallMessages []string `json:"post_install_messages"`
}
type Ubuntu struct {
	Apt   Apt     `json:"apt"`
	Files []File  `json:"files"`
}
type Unsupported struct {
	InstallSteps []InstallStep
	UpdateSteps  []UpdateSteps
}
type Platforms struct {
	Macos   Macos  `json:"macos"`
	Ubuntu  Ubuntu `json:"ubuntu"`
}

type InstallOptions struct {
	OverwriteFiles bool
	BedrockDir string
}

func (e Extension) Load(pkgm string) Extension {
	e.BasePath = filepath.Join(helpers.Home, ".bedrock", "extensions", e.Name)
	extension := LoadManifest(e, pkgm)

	return extension
}

func LoadManifest(e Extension, pkgm string) Extension {
	manifestJson, _ := ioutil.ReadFile(filepath.Join(e.BasePath, "manifest.json"))

	var manifest Manifest
	json.Unmarshal(manifestJson, &manifest)

	switch helpers.CurrentPlatform() {
	case "macos":
		if pkgm == "homebrew" {
			e.InstallSteps = manifest.Platforms.Macos.Homebrew.InstallSteps
			e.PostInstallMessages = manifest.Platforms.Macos.Homebrew.PostInstallMessages
		}
		e.Files = manifest.Platforms.Macos.Files
	case "ubuntu":
		if pkgm == "apt" {
			e.InstallSteps = manifest.Platforms.Ubuntu.Apt.InstallSteps
			e.PostInstallMessages = manifest.Platforms.Ubuntu.Apt.PostInstallMessages
		}
		e.Files = manifest.Platforms.Ubuntu.Files
	}

	return e
}

func (e Extension) Install(options InstallOptions) {
	installResult := e.RunInstallSteps(options)
	syncResult := e.SyncFiles(options)

	if installResult && syncResult {
		if len(e.PostInstallMessages) > 0 {
			message := `
    =========================
    Post-install instructions
    =========================
`
			fmt.Println(helpers.ColorYellow + message + helpers.ColorReset)
			for _, line := range e.PostInstallMessages {
				fmt.Printf("    %s\n", helpers.ColorYellow + line + helpers.ColorReset)
			}
		}
		fmt.Println(e.Name, "-", helpers.ColorGreen + "succeeded" + helpers.ColorReset)
	} else {
		fmt.Println(e.Name, "-", helpers.ColorRed + "failed" + helpers.ColorReset)
	}
}

func (e Extension) RunInstallSteps(options InstallOptions) bool {
	if len(e.InstallSteps) == 0 {
		return true
	}

	fmt.Println(e.Name, "-", helpers.ColorYellow + "installing" + helpers.ColorReset)

	for _, step := range e.InstallSteps {
		pathExpansions := []string{"~", helpers.Home, "$HOME", helpers.Home, "$BEDROCK_DIR", options.BedrockDir}
		command := helpers.ExpandPath(step.Command, pathExpansions...)
		runIf := helpers.ExpandPath(step.RunIf, pathExpansions...)

		fmt.Printf("  %s %s %s\n", "Executing",  helpers.ColorYellow + step.Binary,
			command + helpers.ColorReset)

		if len(runIf) > 0 {
			if _, ifCheckErr := ExecuteRunIfCheck(runIf); ifCheckErr != nil {
				fmt.Printf("    %s\n", helpers.ColorCyan + "Skipping due to runif check" + helpers.ColorReset)
				continue
			}
		}

		// FIXME: the command argument splitting in helpers.ExecuteCommand messes with the natural quoting users would
		//        supply in `-c` argument when setting the binary to something like `sh`.
		out, err := helpers.ExecuteCommand(step.Binary, command)

		var color string
		if err != nil {
			color = helpers.ColorRed
		} else {
			color = helpers.ColorGreen
		}

		if len(out) > 0 {
			for _, line := range strings.Split(string(out), "\n") {
				fmt.Printf("    %s\n", color+line+helpers.ColorReset)
			}
		}
	}

	return true
}

func ExecuteRunIfCheck(command string) (string, error) {
	out, err := exec.Command("sh", "-c", command).CombinedOutput()

	return string(out), err
}

func (e Extension) SyncFiles(options InstallOptions) bool {
	if len(e.Files) > 0 {
		fmt.Println("  Syncing files")
	}

	pathExpansions := []string{"~", helpers.Home, "$HOME", helpers.Home, "$BEDROCK_DIR", options.BedrockDir}

	// TODO: Support skipping all.
	// TODO: Support skipping all for the current extension.
	// TODO: Support overwriting all for the current extension.
	for _, f := range e.Files {
		var source string

		if f.Operation == "remote" {
			source = f.Source
		} else {
			source = filepath.Join(e.BasePath,  helpers.ExpandPath(f.Source, pathExpansions...))
			if ! helpers.Exists(source) {
				fmt.Println("    " + helpers.ColorRed + source + " does not exist, skipping." + helpers.ColorReset)
				return false
			}
		}

		destination := helpers.ExpandPath(f.Target, pathExpansions...)
		destinationExists := helpers.Exists(destination)

		if destinationExists && !options.OverwriteFiles {
			fmt.Printf("    %s already exists. Attempt to overwrite? y/n%s ", helpers.ColorYellow + destination,
				helpers.ColorReset)
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			fmt.Println("")
			if strings.TrimSpace(response) != "y" {
				fmt.Println("    " + helpers.ColorCyan + "Skipping" + " " + destination + helpers.ColorReset)
				break
			}
		}

		destinationBasePath := filepath.Dir(destination)
		if ! helpers.Exists(destinationBasePath) {
			os.MkdirAll(destinationBasePath, os.FileMode(0744))
		}

		os.Remove(destination)

		// FIXME: there's no guard against no operation being specified in the manifest
		switch f.Operation {
		case "copy":
			helpers.Copy(source, destination)
		case "symlink":
			os.Symlink(source, destination)
		case "remote":
			if err := helpers.Download(source, destination); err != nil {
				fmt.Printf("    %s%s %s - %v%s\n",
					helpers.ColorRed,
					"Unable to download",
					source,
					err,
					helpers.ColorReset)
				return false
			}
		}

		fmt.Printf("    %s %s\n", helpers.ColorYellow + f.Operation,
			f.Source + " -> "+ f.Target + helpers.ColorReset)
	}

	return true
}
