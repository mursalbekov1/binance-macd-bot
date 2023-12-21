package main

import (
	binance2 "binance_tg/binance"
	"fmt"
	"log"
)

func main() {

	//https://api.telegram.org/bot<token>/METHOD_NAME
	botToken := "6447935217:AAFVtGASsalRTgVXXUae31OrrIYXGyYnBbQ"
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken
	offset := 0
	isRunning := false

	for {
		updates, err := binance2.GetUpdates(botUrl, offset)
		if err != nil {
			log.Println("Something went wrong: ", err.Error())
		}
		if !isRunning {
			binance2.CheckState(botUrl)
			log.Println("проверено")
			isRunning = true
		}
		for _, update := range updates {
			err = binance2.Respond(botUrl, update)
			offset = update.UpdateId + 1
		}
		fmt.Println(updates)
	}
}
