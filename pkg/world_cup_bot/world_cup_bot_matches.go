package world_cup_bot

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

type Match struct {
	Venue           string `json:"venue"`
	Location        string `json:"location"`
	Status          string `json:"status"`
	Time            string `json:"time"`
	FifaID          string `json:"fifa_id"`
	Datetime        string `json:"datetime"`
	LastEventUpdate string `json:"last_event_update_at"`
	LastScoreUpdate string `json:"last_score_update_at"`
	HomeTeam        Team   `json:"home_team"`
	AwayTeam        Team   `json:"away_team"`
	Winner          string `json:"winner"`
	WinnerCode      string `json:"winner_code"`
}

type Team struct {
	Country     string `json:"country"`
	CountryCode string `json:"code"`
	Goals       int    `json:"goals"`
}

func todaysMatches(update tgbotapi.Update, bot *tgbotapi.BotAPI, redisClient *redis.Client) {
	getMatches(update, bot, redisClient, "today_matches")
}

func currentMatches(update tgbotapi.Update, bot *tgbotapi.BotAPI, redisClient *redis.Client) {
	getMatches(update, bot, redisClient, "current_matches")
}

func getMatches(update tgbotapi.Update, bot *tgbotapi.BotAPI, redisClient *redis.Client, getType string) {
	matchesStr, err := redisClient.Get(getType).Result()
	if err != nil {
		log.Printf("Error: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "whoops, something went terribly wrong")
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
		return
	}
	matches := make([]Match, 0)
	json.Unmarshal([]byte(matchesStr), &matches)
	fmt.Println("Matches", matches)
	result := ""
	for _, match := range matches {
		result += fmt.Sprintf(`%s (%d) - (%d) %s, %s
Match time: %s

`, match.HomeTeam.Country, match.HomeTeam.Goals, match.AwayTeam.Goals, match.AwayTeam.Country, match.Time, match.Datetime)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, result)
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}
