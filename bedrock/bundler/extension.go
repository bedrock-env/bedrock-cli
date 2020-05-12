package bundler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bedrock-env/bedrock-cli/bedrock/helpers"
)

type Extension struct {
	Name string
	Git string
	Branch string
	Ref string
	Tag string
	Path string
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
	InstallSteps []InstallStep   `json:"install_steps"`
	UpdateSteps  []UpdateSteps   `json:"update_steps"`
	PostInstallMessages []string `json:"post_install_messages"`
}
type File struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Operation string `json:"operation"`
}
type Macos struct {
	Homebrew Homebrew `json:"homebrew"`
	Files    []File   `json:"files"`
}
type Apt struct {
	InstallSteps []InstallStep   `json:"install_steps"`
	UpdateSteps  []UpdateSteps   `json:"update_steps"`
	PostInstallMessages []string `json:"post_install_messages"`
}
type Ubuntu struct {
	Apt   Apt    `json:"apt"`
	Files []File `json:"files"`
}
type Unsupported struct {
	InstallSteps []InstallStep
	UpdateSteps  []UpdateSteps
}
type Platforms struct {
	Macos  Macos  `json:"macos"`
	Ubuntu Ubuntu `json:"ubuntu"`
}

func (e Extension) Init(options Options) Extension {
	e.getSource(options)
	bundlePath := filepath.Join(options.BedrockDir, "bundle")

	if len(e.Path) > 0 {
		e.BasePath = helpers.ExpandPath(e.Path)
	} else {
		e.BasePath = filepath.Join(bundlePath, e.Name)
	}

	extension := e.loadManifest(options.PackageManager)

	return extension
}

func (e Extension) Install(options Options) bool {
	installResult := e.runInstallSteps(options)
	syncResult := e.syncFiles(options)

	if installResult && syncResult {
		if len(e.PostInstallMessages) > 0 {
			message := `
    =========================
    Post-install instructions
    =========================
`
			fmt.Println(helpers.ColorYellow + message + helpers.ColorReset)
			for _, line := range e.PostInstallMessages {
				fmt.Printf("    %s\n", helpers.ColorYellow+ line +helpers.ColorReset)
			}
		}
		fmt.Println(e.Name, "-", helpers.ColorGreen+ "succeeded" +helpers.ColorReset)

		return true
	} else {
		fmt.Println(e.Name, "-", helpers.ColorRed+ "failed" +helpers.ColorReset)
	}

	return false
}

func (e Extension) loadManifest(pkgm string) Extension {
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

func (e Extension) getSource(options Options) {
	if len(e.Path) > 0 {
		return
	}

	if e.Git == "" {
		e.Git = fmt.Sprintf("https://github.com/bedrock-env/%s.git", e.Name)
	}

	command := fmt.Sprintf("git -C %s clone %s %s", filepath.Join(options.BedrockDir, "bundle"), e.Git, e.Name)
	var checkoutTarget string

	switch {
	case e.Branch != "":
		checkoutTarget = e.Branch
	case e.Ref != "":
		checkoutTarget = e.Ref
	case e.Tag != "":
		checkoutTarget = e.Tag
	}

	if len(checkoutTarget) > 0  {
		command = fmt.Sprintf("%s && git -C %s checkout %s",
			command,filepath.Join(options.BedrockDir, "bundle", e.Name), checkoutTarget)
	}

	fmt.Println(command)
	out, err := helpers.ExecuteCommandInShell(exec.Command, "zsh", command)
	fmt.Println("out", out)
	fmt.Println("err", err)
}

func (e Extension) runInstallSteps(options Options) bool {
	if len(e.InstallSteps) == 0 {
		return true
	}

	fmt.Println(e.Name, "-", helpers.ColorYellow+ "installing" +helpers.ColorReset)

	for _, step := range e.InstallSteps {
		pathExpansions := []string{"~", helpers.Home, "$HOME", helpers.Home, "$BEDROCK_DIR", options.BedrockDir}
		command := helpers.ExpandPath(step.Command, pathExpansions...)
		runIf := helpers.ExpandPath(step.RunIf, pathExpansions...)

		fmt.Printf("  %s %s %s\n", "Executing",  helpers.ColorYellow+ step.Binary,
			command +helpers.ColorReset)

		if len(runIf) > 0 {
			if _, ifCheckErr := executeRunIfCheck(runIf); ifCheckErr != nil {
				fmt.Printf("    %s\n", helpers.ColorCyan+ "Skipping due to runif check" +helpers.ColorReset)
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

func executeRunIfCheck(command string) (string, error) {
	out, err := exec.Command("sh", "-c", command).CombinedOutput()

	return string(out), err
}

func (e Extension) syncFiles(options Options) bool {
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
			fmt.Printf("    %s already exists. Attempt to overwrite? y/n%s ", helpers.ColorYellow+ destination,
				helpers.ColorReset)
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			fmt.Println("")
			if strings.TrimSpace(response) != "y" {
				fmt.Println("    " + helpers.ColorCyan + "Skipping" + " " + destination + helpers.ColorReset)
				continue
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

		fmt.Printf("    %s %s\n", helpers.ColorYellow+ f.Operation,
			f.Source + " -> "+ f.Target +helpers.ColorReset)
	}

	return true
}
