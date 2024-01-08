package main

import (
	binance2 "binance_tg/binance"
	"binance_tg/logging"
)

func main() {
	botToken := "6447935217:AAFVtGASsalRTgVXXUae31OrrIYXGyYnBbQ"
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken
	offset := 0
	isRunning := false

	uid := logging.GenerateString(32)
	logger, file := logging.CustomLog("", uid)
	defer file.Close()

	for {
		logger.Printf("main: function started\n")
		updates, err := binance2.GetUpdates(botUrl, offset)
		if err != nil {
			logger.Printf("Something went wrong: %s\n", err.Error())
		}
		if !isRunning {
			binance2.CheckState(botUrl, uid)
			logger.Printf("Checked\n")
			isRunning = true
		}
		logger.Printf("processing updates started\n")
		for _, update := range updates {
			err = binance2.Respond(botUrl, update, uid)
			offset = update.UpdateId + 1
		}
		logger.Printf("processing updates complete\nd")
	}
}
