package api

import (
	"binance_tg/internal/functions"
	"context"
	"github.com/adshao/go-binance/v2"
	"log"
	"strconv"
	"time"
)

func GetMACD(client *binance.Client, symbol string, interval string, limit int) float64 {

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
		macd, signalLine := functions.CalculateMACD(prices, shortPeriod, longPeriod, signalPeriod)

		time.Sleep(time.Second)
		return macd[len(macd)-1] - signalLine[len(signalLine)-1]
	}

}
