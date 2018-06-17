package main

import (
	"log"
	"os"
	"path"

	"github.com/amitizle/telegram-world-cup-bot/pkg/world_cup_api"
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
	Run:   run,
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
	viper.SetDefault("port", 9988)
	viper.SetDefault("host", "0.0.0.0")
	viper.SetDefault("telegram_token", "")
	viper.SetDefault("redis_host", "127.0.0.1")
	viper.SetDefault("redis_port", 6379)
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

func run(cmd *cobra.Command, args []string) {
	runPoller(cmd, args)
	runBot(cmd, args)
}

func runPoller(cmd *cobra.Command, args []string) {
	redisHost := viper.GetString("redis_host")
	redisPort := viper.GetInt("redis_port")
	world_cup_api.StartPolling(redisHost, redisPort)
}

func runBot(cmd *cobra.Command, args []string) {
	token := viper.GetString("telegram_token")
	redisHost := viper.GetString("redis_host")
	redisPort := viper.GetInt("redis_port")
	localHost := viper.GetString("host")
	localPort := viper.GetInt("port")
	webhookAddr := viper.GetString("webhook_address")
	err := world_cup_bot.Start(webhookAddr, localHost, localPort, token, redisHost, redisPort)
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
