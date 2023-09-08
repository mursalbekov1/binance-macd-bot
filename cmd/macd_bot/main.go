package main

import (
	"binance_tg/internal/api"
	"fmt"
	"log"
)

func main() {
	botToken := "6447935217:AAFVtGASsalRTgVXXUae31OrrIYXGyYnBbQ"
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken

	for {
		updates, err := api.GetUpdates(botUrl)
		if err != nil {
			log.Println("Something went wrong: ", err.Error())
		}
		fmt.Println(updates)
	}
}
