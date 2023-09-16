package main

import (
	"binance_tg/internal/api"
	"fmt"
	"log"
)

// binance api

func main() {

	//https://api.telegram.org/bot<token>/METHOD_NAME
	botToken := "6447935217:AAFVtGASsalRTgVXXUae31OrrIYXGyYnBbQ"
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken
	offset := 0

	for {
		updates, err := api.GetUpdates(botUrl, offset)
		if err != nil {
			log.Println("Something went wrong: ", err.Error())
		}
		for _, update := range updates {
			err = api.Respond(botUrl, update)
			offset = update.UpdateId + 1
		}
		fmt.Println(updates)
	}
}
