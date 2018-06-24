package world_cup_bot

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"time"
)

type Event struct {
	ReceiverChatId int64  `json:"receiver_chat_id"`
	Message        string `json:"message"`
}

// Small dumb goroutine that listens to events using Redis PubSub and
// sends those to the given chat id.
// This allows writing micro services pushing updates to the bot without
// deploying new bot version
func subscribeToEvents(bot *tgbotapi.BotAPI, redisClient *redis.Client, channel string) {
	pubsub := redisClient.Subscribe(channel)
	go listen(bot, pubsub)
}

func listen(bot *tgbotapi.BotAPI, pubsub *redis.PubSub) {
	defer pubsub.Close()
	subscr, err := pubsub.ReceiveTimeout(5 * time.Second)
	if err != nil {
		log.Printf("pubsub error while waiting for subscription to be created: %v", err)
		time.Sleep(5 * time.Second)
		listen(bot, pubsub)
	}
	log.Println(subscr)
	log.Printf("Listening to Redis events")
	for {
		redisMessage, err := pubsub.ReceiveMessage()
		if err != nil {
			log.Printf("pubsub error: %v", err)
			continue
		}

		var event Event
		json.Unmarshal([]byte(redisMessage.Payload), &event)
		log.Printf("received message from Redis PubSub, sends the message to chat %d", event.ReceiverChatId)
		msg := tgbotapi.NewMessage(event.ReceiverChatId, event.Message)
		bot.Send(msg)
	}
}
