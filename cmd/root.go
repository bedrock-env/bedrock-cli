package cmd

import (
	"fmt"
	"github.com/bedrock-env/bedrock-cli/helpers"
	"github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const configFileName string = ".bedrock.json"
const minZSHVersion = "5.0"

var packageManager string
var BedrockDir string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bedrock-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {
	//	// viper.Set("param1", "value1")
	//	fmt.Println("Pkg Manager:", packageManager)
	//	fmt.Println("author:", author)
	//	//package_managers.InstallPackage("foo")
	//	// package_manager_provider.InstallPackages("foo")
	//},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var configPath = filepath.Join(helpers.Home, configFileName)

	if firstRun(configPath) {
		fmt.Println("It looks like this might be the first time Bedrock has been run.")
		fmt.Print("Checking Bedrock requirements...\n\n")
		if !meetRequirements() {
			fmt.Printf("%sRequirements not satisfied. Exiting.%s\n", helpers.ColorRed, helpers.ColorReset)
			os.Exit(1)
		}
	}

	viper.WriteConfigAs(configPath)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&packageManager, "package-manager", helpers.DefaultPkgManager(), "Desired package manager")
	viper.BindPFlag("package-manager", rootCmd.PersistentFlags().Lookup("package-manager"))

	rootCmd.PersistentFlags().StringVar(&BedrockDir, "bedrockdir", filepath.Join(helpers.Home, ".bedrock"), "The Bedrock base directory")
	viper.BindPFlag("bedrockdir", rootCmd.PersistentFlags().Lookup("bedrockdir"))
}

func initConfig() {
	configFilePath := filepath.Join(helpers.Home, configFileName)

	viper.SetConfigFile(configFilePath)
	viper.AutomaticEnv()
}

func meetRequirements() bool {
	zshCheckResult := zshDetected()

	return zshCheckResult
}

func zshDetected() bool {
	detected := false
	result, err := helpers.ExecuteInShell("zsh", "echo $ZSH_VERSION")

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

//func checkPackageManager() {
//
//}

func firstRun(configPath string) bool {
	return !helpers.Exists(configPath)
}