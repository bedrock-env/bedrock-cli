package bedrock

import (
	"fmt"
	"github.com/bedrock-env/bedrock-cli/bedrock/helpers"
	"github.com/hashicorp/go-version"
	"github.com/spf13/viper"
	"time"
)

const VERSION = "0.0.1-alpha.4"

func CheckCLI() bool {
	updateAvailable := false
	currentTime := time.Now().UTC()
	lastUpdateCheckAt := viper.GetString("settings.last_cli_update_check_at")

	if len(lastUpdateCheckAt) == 0 {
		updateAvailable = checkCLIHasUpdate()
	} else {
		lastUpdateCheckTime, timeParseErr := time.Parse(time.RFC3339, lastUpdateCheckAt)

		if timeParseErr == nil {
			diff := currentTime.Sub(lastUpdateCheckTime)
			if diff.Hours() > cliUpdateCheckInterval {
				updateAvailable = checkCLIHasUpdate()
			}
		}
	}

	return updateAvailable
}

func checkCLIHasUpdate() bool {
	tag, err := getLatestRelease(cliLatestReleaseURL)
	if err != nil {
		fmt.Println("Bedrock CLI update check failed.")
		fmt.Println(err)

		return false
	}
	latestVersion, _ := version.NewVersion(tag)
	currentVersion, _ := version.NewVersion(VERSION)

	if latestVersion.LessThanOrEqual(currentVersion) {
		viper.Set("settings.last_cli_update_check_at", time.Now().UTC())
		viper.WriteConfig()
		helpers.RewriteConfig()

		return false
	}

	return true
}
