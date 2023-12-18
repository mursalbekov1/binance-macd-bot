package binance

import (
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
		macd, signalLine := CalculateMACD(prices, shortPeriod, longPeriod, signalPeriod)

		time.Sleep(time.Minute)
		return macd[len(macd)-1] - signalLine[len(signalLine)-1]
	}

}

func CalculateEMA(data []float64, period int) []float64 {
	ema := make([]float64, len(data))
	multiplier := 2.0 / float64(period+1)

	// Вычисляем начальное EMA как простое скользящее среднее
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += data[i]
	}
	ema[period-1] = sum / float64(period)

	// Вычисляем EMA для остальных элементов
	for i := period; i < len(data); i++ {
		ema[i] = (data[i]-ema[i-1])*multiplier + ema[i-1]
	}

	return ema
}

func CalculateMACD(data []float64, shortPeriod, longPeriod, signalPeriod int) ([]float64, []float64) {
	shortEMA := CalculateEMA(data, shortPeriod)
	longEMA := CalculateEMA(data, longPeriod)

	// Вычисляем MACD как разницу между short EMA и long EMA
	var macd []float64
	for i := 0; i < len(data); i++ {
		macd = append(macd, shortEMA[i]-longEMA[i])
	}

	// Вычисляем сигнальную линию MACD как EMA от MACD
	signalLine := CalculateEMA(macd, signalPeriod)

	return macd, signalLine
}
