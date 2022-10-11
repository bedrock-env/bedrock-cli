package bundler

import (
	"errors"
	"fmt"
	"github.com/bedrock-env/bedrock-cli/bedrock/helpers"
	yamlv3 "gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Extension struct {
	Name string
	Git  string
	Path string
	//Archive             string
	Branch              string
	Ref                 string
	Tag                 string
	InstallSteps        []InstallStep
	PostInstallMessages []string
	SourcePath          string
}

type ExtensionManifest struct {
	Author struct {
		Name  string
		Email string
	}
	Platforms []string
	Setup     struct {
		Macos []InstallStep
	}
}

type InstallStep struct {
	Name        string
	Command     string
	RunIf       string
	PostMessage string
	Files       []File
}

type File struct {
	Source    string
	Target    string
	Operation string
}

func (e Extension) Validate() []error {
	var validationErrors []error

	if e.Path == "" && e.Git == "" {
		validationErrors = append(validationErrors, errors.New("path or git must must be specified"))
	}

	return validationErrors
}

func (e Extension) Prepare(options Options) (Extension, error) {
	sourceErr := e.getSource()
	if sourceErr != nil {
		return Extension{}, sourceErr
	}

	bundlePath := filepath.Join(options.BedrockDir, "bundlenew")

	if len(e.Path) > 0 {
		e.SourcePath = helpers.ExpandPath(e.Path)
	} else {
		e.SourcePath = filepath.Join(bundlePath, e.Name)
	}

	extension := e.hydrate()

	return extension, nil
}

func (e Extension) Setup(options Options) bool {
	for _, step := range e.InstallSteps {
		fmt.Println(helpers.BasicStyle.Render(step.Name))

		command := helpers.ExpandPath(step.Command)
		runIf := helpers.ExpandPath(step.RunIf)

		if len(runIf) > 0 {
			if _, ifCheckErr := executeRunIfCheck(runIf); ifCheckErr != nil {
				fmt.Println(helpers.WarnStyle.Render("Skipping due to runif check"))
				//if len(out) > 0 {
				//	fmt.Println(out)
				//	fmt.Print(ifCheckErr)
				//}

				continue
			}
		}

		// FIXME: the command argument splitting in helpers.ExecuteCommand messes with the natural quoting users would
		//        supply in `-c` argument when setting the binary to something like `sh`.
		out, err := helpers.ExecuteCommandInShell(exec.Command, "zsh", command)

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

	installResult := e.runSteps(options)

	var syncResult bool
	if options.SkipFiles {
		syncResult = true
		fmt.Println("    " + helpers.ColorCyan + "Skipping syncing files" + helpers.ColorReset)
	} else {
		syncResult = e.syncFiles(options)
	}

	if installResult && syncResult {
		if len(e.PostInstallMessages) > 0 {
			message := `
    =========================
    Post-install instructions
    =========================
`
			fmt.Println(helpers.ColorYellow + message + helpers.ColorReset)
			for _, line := range e.PostInstallMessages {
				fmt.Printf("    %s\n", helpers.ColorYellow+line+helpers.ColorReset)
			}
		}
		fmt.Println(e.Name, "-", helpers.ColorGreen+"succeeded"+helpers.ColorReset)

		return true
	} else {
		fmt.Println(e.Name, "-", helpers.ColorRed+"failed"+helpers.ColorReset)
	}

	return false
}

func (e Extension) hydrate() Extension {
	// TODO: err when no manifest is found
	path := filepath.Join(e.SourcePath, "manifest.yaml")
	manifestJson, _ := os.ReadFile(path)

	var manifest ExtensionManifest
	yamlv3.Unmarshal(manifestJson, &manifest)

	switch helpers.CurrentPlatform() {
	case "macos":
		e.InstallSteps = manifest.Setup.Macos
	}

	return e
}

func (e Extension) getSource() error {
	var err error

	switch {
	case e.Path != "":
		return nil
	case e.Git != "":
		err = e.getSourceFromGit()
		//case e.Archive != "":
		//e.getSourceFromArchive()
	}

	return err
}

func (e Extension) getSourceFromGit() error {
	command := fmt.Sprintf("git -C %s clone %s %s", filepath.Join(helpers.BedrockDir, "bundle"), e.Git, e.Name)
	var checkoutTarget string

	switch {
	case e.Branch != "":
		checkoutTarget = e.Branch
	case e.Ref != "":
		checkoutTarget = e.Ref
	case e.Tag != "":
		checkoutTarget = e.Tag
	}

	if len(checkoutTarget) > 0 {
		command = fmt.Sprintf("%s && git -C %s checkout %s",
			command, filepath.Join(helpers.BedrockDir, "bundle", e.Name), checkoutTarget)
	}

	// TODO: maybe pass in the output pipe and live write the command output
	_, err := helpers.ExecuteCommandInShell(exec.Command, "zsh", command)

	return err
}

func (e Extension) runSteps(options Options) bool {
	//if len(e.InstallSteps) == 0 {
	//	return true
	//}
	//
	//fmt.Println(e.Name, "-", helpers.ColorYellow+"installing"+helpers.ColorReset)
	//
	//for _, step := range e.InstallSteps {
	//	// pathExpansions := []string{"~", helpers.Home, "$HOME", helpers.Home, "$BEDROCK_DIR", options.BedrockDir}
	//	command := helpers.ExpandPath(step.Command)
	//	runIf := helpers.ExpandPath(step.RunIf)
	//
	//	fmt.Printf("  %s %s %s\n", "Executing", helpers.ColorYellow+step.Binary,
	//		command+helpers.ColorReset)
	//
	//	if len(runIf) > 0 {
	//		if out, ifCheckErr := executeRunIfCheck(runIf); ifCheckErr != nil {
	//			fmt.Printf("    %s\n", helpers.ColorCyan+"Skipping due to runif check"+helpers.ColorReset)
	//			if len(out) > 0 {
	//				fmt.Println(out)
	//				fmt.Print(ifCheckErr)
	//			}
	//
	//			continue
	//		}
	//	}
	//
	//	// FIXME: the command argument splitting in helpers.ExecuteCommand messes with the natural quoting users would
	//	//        supply in `-c` argument when setting the binary to something like `sh`.
	//	out, err := helpers.ExecuteCommand(step.Binary, command)
	//
	//	var color string
	//	if err != nil {
	//		color = helpers.ColorRed
	//	} else {
	//		color = helpers.ColorGreen
	//	}
	//
	//	if len(out) > 0 {
	//		for _, line := range strings.Split(string(out), "\n") {
	//			fmt.Printf("    %s\n", color+line+helpers.ColorReset)
	//		}
	//	}
	//}

	return true
}

func executeRunIfCheck(command string) (string, error) {
	out, err := exec.Command("sh", "-c", command).CombinedOutput()

	return string(out), err
}

func (e Extension) syncFiles(options Options) bool {
	//if len(e.Files) > 0 {
	//	fmt.Println("  Syncing files")
	//}
	//
	//pathExpansions := []string{"~", helpers.Home, "$HOME", helpers.Home, "$BEDROCK_DIR", options.BedrockDir}
	//
	//// TODO: Support overwriting all for the current extension.
	//skipAll := false
	//overwriteAll := false
	//
	//for _, f := range e.Files {
	//	if skipAll {
	//		continue
	//	}
	//
	//	var source string
	//
	//	if f.Operation == "remote" {
	//		source = f.Source
	//	} else {
	//		source = filepath.Join(e.BasePath, helpers.ExpandPath(f.Source, pathExpansions...))
	//		if !helpers.Exists(source) {
	//			fmt.Println("    " + helpers.ColorRed + source + " does not exist, skipping." + helpers.ColorReset)
	//			return false
	//		}
	//	}
	//
	//	destination := helpers.ExpandPath(f.Target, pathExpansions...)
	//	destinationExists := helpers.Exists(destination)
	//
	//	if !options.OverwriteFiles && !overwriteAll && destinationExists {
	//		fmt.Printf("    %s already exists. Attempt to overwrite? y/n/(s)kip remaining/(O)verwrite remaining)%s ", helpers.ColorYellow+destination,
	//			helpers.ColorReset)
	//		reader := bufio.NewReader(os.Stdin)
	//		response, _ := reader.ReadString('\n')
	//		response = strings.TrimSpace(response)
	//		fmt.Println("")
	//		if response == "s" {
	//			skipAll = true
	//			fmt.Println("    " + helpers.ColorCyan + "Skipping remaining files for extension\n" + helpers.ColorReset)
	//			continue
	//		} else if response != "y" {
	//			fmt.Println("    " + helpers.ColorCyan + "Skipping" + " " + destination + "\n" + helpers.ColorReset)
	//			continue
	//		} else if response != "O" {
	//			overwriteAll = true
	//			fmt.Println("    " + helpers.ColorYellow + "Overwriting all files for extension\n" + helpers.ColorReset)
	//			continue
	//		}
	//	}
	//
	//	destinationBasePath := filepath.Dir(destination)
	//	if !helpers.Exists(destinationBasePath) {
	//		os.MkdirAll(destinationBasePath, os.FileMode(0744))
	//	}
	//
	//	os.Remove(destination)
	//
	//	// FIXME: there's no guard against no operation being specified in the manifest
	//	switch f.Operation {
	//	case "copy":
	//		helpers.Copy(source, destination)
	//	case "symlink":
	//		os.Symlink(source, destination)
	//	case "remote":
	//		if err := helpers.Download(source, destination); err != nil {
	//			fmt.Printf("    %s%s %s - %v%s\n",
	//				helpers.ColorRed,
	//				"Unable to download",
	//				source,
	//				err,
	//				helpers.ColorReset)
	//			return false
	//		}
	//	}
	//
	//	fmt.Printf("    %s %s\n", helpers.ColorYellow+f.Operation,
	//		f.Source+" -> "+f.Target+helpers.ColorReset)
	//}

	return true
}
