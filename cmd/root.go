package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bedrock-env/bedrock-cli/bedrock/helpers"
)

var rootCmd = &cobra.Command{
	Use: "bedrock",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filepath.Join(helpers.Home, ".config", "bedrock", "config.yaml"))
	viper.AutomaticEnv()
	viper.ReadInConfig()
}
