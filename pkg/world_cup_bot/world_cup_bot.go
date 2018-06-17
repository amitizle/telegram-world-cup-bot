package world_cup_bot

import (
	"errors"
	"fmt"
	"github.com/amitizle/telegram-world-cup-bot/internal/http_client"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	// "net/http"
)

func Start(host string, port int, telegramToken string) error {
	if telegramToken == "" {
		return errors.New("Bot token is missing")
	}
	httpClient, err := world_cup_http_client.New("")
	if err != nil {
		return err
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

	handleUpdates(updates, bot, httpClient)

	return nil
}

func handleUpdates(updateChannel tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI, httpClient *world_cup_http_client.HTTPClient) {
	for update := range updateChannel {
		handleUpdate(update, bot, httpClient)
	}
}

func handleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI, httpClient *world_cup_http_client.HTTPClient) {
	fmt.Println("Message", update.Message)
	if update.Message == nil {
		return
	}
	fmt.Printf("IkoPico")
	fmt.Printf("Command", update.Message.Command())
	switch update.Message.Command() {
	case "today":
		todaysMatches(update, bot, httpClient)
	case "version":
		botVersion(update, bot)
	default:
		log.Printf("No handler for %v", update.Message)
	}
}
