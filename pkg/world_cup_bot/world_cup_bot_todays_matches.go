package world_cup_bot

import (
	"encoding/json"
	"fmt"
	"github.com/amitizle/telegram-world-cup-bot/internal/http_client"
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

func todaysMatches(update tgbotapi.Update, bot *tgbotapi.BotAPI, httpClient *world_cup_http_client.HTTPClient) {
	response, err := httpClient.Get("/matches/today", map[string]string{})
	if err != nil {
		log.Printf("Error: %v", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "whoops, something went terribly wrong")
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
		return
	}
	matches := make([]Match, 0)
	json.Unmarshal(response.Body, &matches)
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
