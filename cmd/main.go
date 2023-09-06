package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

var gBot *tgbotapi.BotAPI
var gToken string
var gChatID int64

func isStartMessage(update *tgbotapi.Update) bool {
	return update.Message != nil && update.Message.Text == "/start"
}

func init() {
	_ = os.Setenv("binance_macd_indicator_bot", "6447935217:AAFVtGASsalRTgVXXUae31OrrIYXGyYnBbQ")

	gToken = os.Getenv("binance_macd_indicator_bot")

	var err error
	if gBot, err = tgbotapi.NewBotAPI(gToken); err != nil {
		log.Panic(err)
	}
	gBot.Debug = true
}

func main() {
	var err error
	gBot, err := tgbotapi.NewBotAPI(gToken)
	if err != nil {
		log.Panic(err)
	}

	gBot.Debug = true

	log.Printf("Authorized on account %s", gBot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	for update := range gBot.GetUpdatesChan(u) {
		if isStartMessage(&update) { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			gBot.Send(msg)
		}
	}
}
