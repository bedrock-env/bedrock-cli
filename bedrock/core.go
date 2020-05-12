package bedrock

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-version"
	"github.com/spf13/viper"

	"github.com/bedrock-env/bedrock-cli/bedrock/helpers"
)

const minZSHVersion = "5.0"
const bedrockMinVersion = "0.0.1-alpha"

func CheckFirstRun(configPath string) {
	if !helpers.Exists(configPath) {
		fmt.Println("It looks like this might be the first time Bedrock has been run.")
		fmt.Print("Checking Bedrock requirements...\n\n")
		if !meetRequirements() {
			fmt.Printf("%sRequirements not satisfied. Exiting.%s\n", helpers.ColorRed, helpers.ColorReset)
			os.Exit(1)
		}
	}

	viper.WriteConfigAs(configPath)
}

func meetRequirements() bool {
	zshCheckResult := zshDetected()

	return zshCheckResult
}

func zshDetected() bool {
	detected := false
	result, err := helpers.ExecuteCommandInShell(exec.Command, "zsh", "echo $ZSH_VERSION")

	if err == nil {
		zshVersion, _ := version.NewVersion(result)
		requiredVersion, _ := version.NewVersion(minZSHVersion)

		if zshVersion.GreaterThanOrEqual(requiredVersion) {
			fmt.Printf("%s\u2714%s ZSH %s detected\n", helpers.ColorGreen, helpers.ColorReset, zshVersion)
			return true
		}
	}

	fmt.Printf("%s\u0078%s ZSH %s was not detected\n", helpers.ColorRed, helpers.ColorReset, minZSHVersion)

	return detected
}
