package world_cup_bot

import (
	"gopkg.in/telegram-bot-api.v4"
)

var (
	version = "0.3.1"
)

func botVersion(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, version)
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}
