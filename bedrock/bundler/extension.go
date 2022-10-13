package bundler

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/bedrock-env/bedrock-cli/bedrock/helpers"
	"github.com/charmbracelet/lipgloss"
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
	Branch       string
	Ref          string
	Tag          string
	InstallSteps []InstallStep
	InstallNotes string
	SourcePath   string
}

type ExtensionManifest struct {
	Author struct {
		Name  string
		Email string
	}
	Platforms []string
	Setup     struct {
		Macos struct {
			InstallNotes string `yaml:"install_notes"`
			Steps        []InstallStep
		}
	}
}

type InstallStep struct {
	Name    string
	Command string
	RunIf   string
	Files   []File
}

type File struct {
	Source    string
	Target    string
	Operation string
}

func (e *Extension) Validate() []error {
	var validationErrors []error

	if e.Path == "" && e.Git == "" {
		validationErrors = append(validationErrors, errors.New("path or git must must be specified"))
	}

	return validationErrors
}

func (e *Extension) Prepare(options Options) error {
	sourceErr := e.getSource()
	if sourceErr != nil {
		return sourceErr
	}

	bundlePath := filepath.Join(options.BedrockDir, "bundle")

	if len(e.Path) > 0 {
		e.SourcePath = helpers.ExpandPath(e.Path)
	} else {
		e.SourcePath = filepath.Join(bundlePath, e.Name)
	}

	e.hydrate()

	return nil
}

func (e *Extension) Setup(options Options) bool {
	succeeded := true

	for _, step := range e.InstallSteps {
		fmt.Println(helpers.ExtensionInstallStep.Render(step.Name))

		command := helpers.ExpandPath(step.Command)
		runIf := helpers.ExpandPath(step.RunIf)

		if len(runIf) > 0 {
			if _, ifCheckErr := executeRunIfCheck(runIf); ifCheckErr != nil {
				fmt.Println(helpers.WarnStyle.MarginLeft(2).Render("Skipping due to runif check"))

				continue
			}
		}

		out, err := helpers.ExecuteCommandInShell(exec.Command, "zsh", command)

		if err != nil {
			succeeded = false

			fmt.Println(helpers.ErrorStyle.MarginLeft(2).Render("Failed!"))
			fmt.Println(helpers.ErrorStyle.MarginLeft(4).Render("- " + out))

			break

		} else {
			if len(out) > 0 {
				fmt.Println(helpers.BasicStyle.MarginLeft(2).Render("- " + out))
			}
		}

		syncSucceeded := syncFiles(step, e.SourcePath, options)

		if !syncSucceeded {
			succeeded = false

			break
		}
	}

	if len(e.InstallNotes) > 0 {
		fmt.Println(lipgloss.NewStyle().MarginLeft(2).Foreground(helpers.COLORWARN).Render(
			fmt.Sprintf("\n%s\n\n%s\n\n%s",
				"=========== Install Notes ===========",
				e.InstallNotes,
				"=====================================",
			)),
		)
	}

	return succeeded
}

func (e *Extension) hydrate() {
	// TODO: err when no manifest is found
	path := filepath.Join(e.SourcePath, "manifest.yaml")
	manifestJson, _ := os.ReadFile(path)

	var manifest ExtensionManifest
	yamlv3.Unmarshal(manifestJson, &manifest)

	switch helpers.CurrentPlatform() {
	case "macos":
		e.InstallSteps = manifest.Setup.Macos.Steps
		e.InstallNotes = manifest.Setup.Macos.InstallNotes
	}
}

func (e *Extension) getSource() error {
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

func (e *Extension) getSourceFromGit() error {
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

func executeRunIfCheck(command string) (string, error) {
	out, err := exec.Command("zsh", "-c", command).CombinedOutput()

	return string(out), err
}

func syncFiles(step InstallStep, sourcePath string, options Options) bool {
	if len(step.Files) == 0 {
		return true
	}

	pathExpansions := []string{"~", helpers.Home, "$HOME", helpers.Home, "$BEDROCK_DIR", options.BedrockDir}

	// TODO: Support overwriting all for the current extension.
	skipAll := false
	overwriteAll := false

	for _, f := range step.Files {
		if skipAll {
			continue
		}

		var source string

		if f.Operation == "remote" {
			source = f.Source
		} else {
			source = filepath.Join(sourcePath, helpers.ExpandPath(f.Source, pathExpansions...))
			if !helpers.Exists(source) {
				fmt.Println("    " + helpers.ColorRed + source + " does not exist, skipping." + helpers.ColorReset)
				return false
			}
		}

		destination := helpers.ExpandPath(f.Target, pathExpansions...)
		destinationExists := helpers.Exists(destination)

		if !options.OverwriteFiles && !overwriteAll && destinationExists {
			fmt.Printf(
				"%s%s",
				lipgloss.NewStyle().MarginLeft(2).Render(fmt.Sprintf("%s exists.", destination)),
				lipgloss.NewStyle().Foreground(helpers.COLORWARN).Bold(true).Render(fmt.Sprintf(" Attempt to overwrite? y/n/(S)kip all/(O)verwrite all) ")))

			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(response)

			if response == "n" {
				fmt.Println(lipgloss.NewStyle().Bold(true).MarginLeft(2).Foreground(helpers.COLORWARN).Render("Skipping " + f.Source))
				continue
			} else if response == "S" {
				skipAll = true
				fmt.Println(lipgloss.NewStyle().Bold(true).MarginLeft(2).Foreground(helpers.COLORWARN).Render("Skipping remaining files"))
				continue
			} else if response == "O" {
				overwriteAll = true
				fmt.Println(lipgloss.NewStyle().Bold(true).MarginLeft(2).Foreground(helpers.COLORWARN).Render("Overwriting all files"))
			}
		}

		destinationBasePath := filepath.Dir(destination)
		if !helpers.Exists(destinationBasePath) {
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
				fmt.Println(lipgloss.NewStyle().Foreground(helpers.COLORSUCCESS).MarginLeft(2).Render(fmt.Sprintf("%s %s - %v\n",
					"Unable to download",
					source,
					err,
				)))
				return false
			}
		}

		fmt.Println(lipgloss.NewStyle().Foreground(helpers.COLORSUCCESS).MarginLeft(2).Render(fmt.Sprintf("%s %s -> %s", f.Operation, f.Source, f.Target)))
	}

	return true
}
