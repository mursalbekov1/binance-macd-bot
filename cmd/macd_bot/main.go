package main

import (
	"binance_tg/internal/functions"
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"log"
	"strconv"
	"time"
)

// binance api
var (
	apiKey    = "DQbMjZ54eTkw7pHIxYYW1UUFWNADxaETaE79C6Ad0VH69QImdQqVLE9rpJ6ZFc28"
	secretKey = "Mdiyex2E6kiQf2jmOSjrXKsTbqTb4SdURt3LqumbqZi3OdNSUhP3I0NTi8dHdBDG"
)

func main() {

	client := binance.NewClient(apiKey, secretKey)
	symbol := "BTCUSDT"
	interval := "1s"
	limit := 100

	for {
		klines, err := client.NewKlinesService().
			Symbol(symbol).
			Interval(interval).
			Limit(limit).
			Do(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		var prices []float64
		for _, kline := range klines {
			closePriceStr := kline.Close
			closePrice, err := strconv.ParseFloat(closePriceStr, 64)
			if err != nil {
				log.Fatal(err)
			}
			prices = append(prices, closePrice)
		}

		// Периоды для вычисления MACD и сигнальной линии
		shortPeriod := 12
		longPeriod := 26
		signalPeriod := 9

		// Вычисляем MACD и сигнальную линию
		macd, _ := functions.CalculateMACD(prices, shortPeriod, longPeriod, signalPeriod)

		signalLine := functions.CalculateSignalLine(macd, signalPeriod)

		// Выводим последний результат MACD
		//fmt.Println("MACD line:", macd[len(macd)-1])
		//fmt.Println("Signal line:", signalLine[len(signalLine)-1])
		fmt.Println("Histogram:", macd[len(macd)-1]-signalLine[len(signalLine)-1])
		fmt.Println("Time: ", time.Now().Format("15:04:05.000"))
		// Задержка в 3 секунды перед следующим запросом
		//time.Sleep(3 * time.Second)
	}

	//https://api.telegram.org/bot<token>/METHOD_NAME
	//botToken := "6447935217:AAFVtGASsalRTgVXXUae31OrrIYXGyYnBbQ"
	//botApi := "https://api.telegram.org/bot"
	//botUrl := botApi + botToken
	//offset := 0
	//
	//for {
	//	updates, err := api.GetUpdates(botUrl, offset)
	//	if err != nil {
	//		log.Println("Something went wrong: ", err.Error())
	//	}
	//	for _, update := range updates {
	//		err = api.Respond(botUrl, update)
	//		offset = update.UpdateId + 1
	//	}
	//	fmt.Println(updates)
	//}
}
