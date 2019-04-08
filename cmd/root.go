package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getEnv() string {
	if os.Getenv("APP_ENV") == "" {
		return "local"
	}
	return os.Getenv("APP_ENV")
}

var cfgPath string

var rootCmd = &cobra.Command{
	Use:   "comic-crawler",
	Short: "Comic crawler",
	Long: "This application is comic crwaler that you can easily know the which comics are updated recently.\n" +
		"And the main reason of the project appear is for me to enter the Go world. Hope you can enjoy it.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")

	// FIXME: https://github.com/spf13/cobra/issues/206
	// rootCmd.MarkPersistentFlagRequired("config")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		fmt.Println("Must to assign a config file path")
		os.Exit(1)
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	// Initialize logger system
	debug := viper.GetBool("system.debug")

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	if getEnv() != "local" {
		// Some log event collection services
	}
	log.SetOutput(os.Stdout)
}
