package bedrock

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/spf13/viper"

	"github.com/bedrock-env/bedrock-cli/bedrock/helpers"
)

const ZshMinVersion = "5.0"
const CoreMinVersion = "0.0.1-alpha"
const coreRepoURL = "https://github.com/bedrock-env/bedrock-core.git"
const coreLatestReleaseURL = "https://api.github.com/repos/bedrock-env/bedrock-core/releases/latest"
const coreUpdateCheckInterval = 168 // 7 days

type CoreCheckResult struct {
	Found bool
	MeetsMinVersion bool
	UpdateAvailable bool
	Version string
}


func Preflight() error {
	coreResult := CheckCore()

	if !coreResult.Found {
		return errors.New("bedrock core not found")
	} else if !coreResult.MeetsMinVersion {
		return errors.New("newer bedrock core required")
	} else if coreResult.UpdateAvailable {
		fmt.Println("A newer version of Bedrock Core is available.")
	}

	return nil
}

func CheckZSH() bool {
	detected := false

	result, err := helpers.ExecuteCommandInShell(exec.Command, "zsh", "echo $ZSH_VERSION")

	if err == nil {
		zshVersion, _ := version.NewVersion(result)
		requiredVersion, _ := version.NewVersion(ZshMinVersion)

		if zshVersion.GreaterThanOrEqual(requiredVersion) {
			fmt.Printf("%s\u2714%s ZSH %s detected\n", helpers.ColorGreen, helpers.ColorReset, zshVersion)
			return true
		}
	}

	fmt.Printf("%s\u0078%s ZSH %s was not detected\n", helpers.ColorRed, helpers.ColorReset, ZshMinVersion)

	return detected
}

func CheckGit() bool {
	detected := false

	version, err := helpers.ExecuteCommandInShell(exec.Command, "zsh", "git version")

	if err != nil {
		fmt.Printf("%s\u0078%s Git was not detected\n", helpers.ColorRed, helpers.ColorReset)
		fmt.Println(err)

		return true
	}

	fmt.Printf("%s\u2714%s %s detected\n", helpers.ColorGreen, helpers.ColorReset, version)

	return detected
}

func CheckCore() CoreCheckResult {
	found := false
	updateAvailable := false
	meetsMinVersion := false
	var ver string

	versionResult, err := ioutil.ReadFile(filepath.Join(helpers.BedrockDir, "VERSION"))

	if err == nil {
		found = true
	}

	bedrockVersion, _ := version.NewVersion(strings.TrimSpace(string(versionResult)))
	requiredVersion, _ := version.NewVersion(CoreMinVersion)

	if bedrockVersion != nil {
		if bedrockVersion.GreaterThanOrEqual(requiredVersion) {
			meetsMinVersion = true
		}
		ver = bedrockVersion.String()
	}

	currentTime := time.Now().UTC()
	lastUpdateCheckAt := viper.GetString("last_update_check_at")
	if len(lastUpdateCheckAt) == 0 {
		updateAvailable = checkCoreHasUpdate(bedrockVersion)
	} else {
		lastUpdateCheckTime, timeParseErr := time.Parse(time.RFC3339, lastUpdateCheckAt)

		if timeParseErr == nil {
			diff := currentTime.Sub(lastUpdateCheckTime)
			if diff.Hours() > coreUpdateCheckInterval {
				updateAvailable = checkCoreHasUpdate(bedrockVersion)
			}
		}
	}

	return CoreCheckResult{
		Found:           found,
		MeetsMinVersion: meetsMinVersion,
		UpdateAvailable: updateAvailable,
		Version:         ver,
	}
}

func InstallCore(interactive bool) (bool, error) {
	if interactive && !promptYN("Install Bedrock Core?") {
		return false, nil
	}

	latestRelease, err := getLatestCoreRelease()
	if err != nil {
		return false, err
	}

	cmd := fmt.Sprintf("git clone %s %s && git -C %s checkout %s", coreRepoURL, helpers.BedrockDir,
		helpers.BedrockDir, latestRelease)
	fmt.Println(cmd)
	_, err = helpers.ExecuteCommandInShell(exec.Command, "zsh", cmd)
	if err != nil {
		return false, err
	}

	setLastCheckAt()

	return true, nil
}

func UpdateCore(interactive bool) (bool, error) {
	if interactive && !promptYN("Update Bedrock Core?") {
		return false, nil
	}

	latestRelease, err := getLatestCoreRelease()
	if err != nil {
		return false, err
	}

	cmd := fmt.Sprintf("cd %s && git clean -f && git fetch --all && git checkout %s", helpers.BedrockDir, latestRelease)
	out, cmdErr := helpers.ExecuteCommandInShell(exec.Command, "zsh", cmd)
	if len(out) > 0 {
		fmt.Println(out)
	}

	if cmdErr != nil {
		return false, err
	}

	return true, nil
}

func checkCoreHasUpdate(v *version.Version) bool {
	tag, err := getLatestCoreRelease()
	if err != nil {
		fmt.Println("Bedrock Core update check failed.")
		fmt.Println(err)
		return false
	}
	latestVersion, _ := version.NewVersion(tag)

	if latestVersion.LessThanOrEqual(v) {
		setLastCheckAt()

		return false
	}

	return true
}

func setLastCheckAt() {
	viper.Set("last_update_check_at", time.Now().UTC())
	viper.WriteConfig()
}

func promptYN(message string) (result bool) {
	fmt.Printf("%s%s (y/n) %s", helpers.ColorYellow, message, helpers.ColorReset)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')

	return strings.TrimSpace(response) == "y"
}

func getLatestCoreRelease() (tag string, error error) {
	resp, err := http.Get(coreLatestReleaseURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var githubData map[string] string
	json.Unmarshal(body, &githubData)

	if err != nil {
		return "", err
	}

	return githubData["tag_name"], nil
}
