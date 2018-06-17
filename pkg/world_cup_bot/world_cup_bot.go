package world_cup_bot

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	// "net/http"
)

func Start(host string, port int, telegramToken string, redisHost string, redisPort int) error {
	if telegramToken == "" {
		return errors.New("Bot token is missing")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", redisHost, redisPort),
		DB:   0,
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatalf("could not connect to Redis: %v", err)
	}
	botAddr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("Starting bot at %s", botAddr)
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		return err
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	handleUpdates(updates, bot, redisClient)

	return nil
}

func handleUpdates(updateChannel tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI, redisClient *redis.Client) {
	for update := range updateChannel {
		handleUpdate(update, bot, redisClient)
	}
}

func handleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI, redisClient *redis.Client) {
	fmt.Println("Message", update.Message)
	if update.Message == nil {
		return
	}
	fmt.Printf("IkoPico")
	fmt.Printf("Command", update.Message.Command())
	switch update.Message.Command() {
	case "today":
		todaysMatches(update, bot, redisClient)
	case "current":
		currentMatches(update, bot, redisClient)
	case "version":
		botVersion(update, bot)
	default:
		log.Printf("No handler for %v", update.Message)
	}
}
