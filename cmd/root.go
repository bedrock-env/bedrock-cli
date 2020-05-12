package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bedrock-env/bedrock-cli/bedrock"
	"github.com/bedrock-env/bedrock-cli/bedrock/helpers"
)

const configFileName string = ".bedrock.json"

var PackageManager string

var rootCmd = &cobra.Command{
	Use:   "bedrock",
}

func Execute() {
	var configPath = filepath.Join(helpers.Home, configFileName)

	bedrock.CheckFirstRun(configPath)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&PackageManager, "package-manager", helpers.DefaultPkgManager(), "Desired package manager")
	viper.BindPFlag("package-manager", rootCmd.PersistentFlags().Lookup("package-manager"))

	// NOTE: Not everything reads from this yet
	rootCmd.PersistentFlags().StringVar(&helpers.BedrockDir, "bedrockdir", filepath.Join(helpers.Home, ".bedrock"), "The Bedrock base directory")
	viper.BindPFlag("bedrockdir", rootCmd.PersistentFlags().Lookup("bedrockdir"))
}

func initConfig() {
	configFilePath := filepath.Join(helpers.Home, configFileName)

	viper.SetConfigFile(configFilePath)
	viper.AutomaticEnv()
}
