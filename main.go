package main

import (
	"log"
	"os"
	"path"

	"github.com/amitizle/telegram-world-cup-bot/pkg/world_cup_bot"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile            string
	defaultConfigPath     = "/etc/default/"
	defaultConfigFilename = "world_cup_bot.yml"
)

var rootCmd = &cobra.Command{
	Use:   "world-cup-bot",
	Short: "Pff world cup bot for Telegram",
	Run:   runBot,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", path.Join(defaultConfigPath, defaultConfigFilename), "full path to config file")
	viper.SetDefault("port", 8080)
	viper.SetDefault("host", "127.0.0.1")
	viper.SetDefault("telegram_token", "")
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath(defaultConfigPath)
		viper.SetConfigName(defaultConfigFilename)
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix("world_cup")

	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}
}

func runBot(cmd *cobra.Command, args []string) {
	token := viper.GetString("telegram_token")
	err := world_cup_bot.Start(viper.GetString("host"), viper.GetInt("port"), token)
	if err != nil {
		log.Fatalf("faild to start bot: %v", err)
		os.Exit(1)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
