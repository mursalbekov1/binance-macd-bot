package binance

import (
	"binance_tg/logging"
	"binance_tg/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetMACD(client *binance.Client, symbol string, interval string, limit int) float64 {

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

	//time.Sleep(5 * time.Minute)
	return macd[len(macd)-1] - signalLine[len(signalLine)-1]

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

// GetMACDLoop Logging done
func GetMACDLoop(botUrl string, chatID int64, uid string) {
	state := getUserState(chatID)
	logger, file := logging.CustomLog(`chatId=`+fmt.Sprint(chatID), uid)
	defer file.Close()

	for state.IsRunning {
		macdValue := GetMACD(client, symbol, interval, limit)

		if (macdValue > 0 && state.PrevMACDValue <= 0) || (macdValue <= 0 && state.PrevMACDValue > 0) {
			logger.Println("Flag processing true")
			var botMessage models.BotMessage
			if macdValue > 0 {
				botMessage = models.BotMessage{
					ChatId: int(chatID),
					Text:   "Значение MACD поднялось на зеленую отметку 🟢 \n" + "Текущее значение: " + strconv.FormatFloat(macdValue, 'f', -1, 64),
				}
				logger.Println("Green flag")
			} else {
				botMessage = models.BotMessage{
					ChatId: int(chatID),
					Text:   "Значение MACD опустилось на красную отметку 🔴 \n" + "Текущее значение: " + strconv.FormatFloat(macdValue, 'f', -1, 64),
				}
				logger.Println("Red flag")
			}
			buf, err := json.Marshal(botMessage)
			if err != nil {
				logger.Println("Ошибка при маршалинге сообщения:", err)
				continue
			}
			_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
			if err != nil {
				logger.Println("Ошибка при отправке сообщения:", err)
			}
			logger.Println("Flag notified successfully")
		}

		setPrevMACDValue(chatID, macdValue)
		logger.Println(`Previous MACD value set - value ` + fmt.Sprint(macdValue))

		time.Sleep(time.Minute * 10)
	}
}

// GetMACDLoopRed logging done
//func GetMACDLoopRed(botUrl string, chatID int64, uid string) {
//	state := getUserState(chatID)
//	logger, file := logging.CustomLog(`chatId=`+fmt.Sprint(chatID), uid)
//	defer file.Close()
//
//	for state.IsRunning {
//		macdValue := GetMACD(client, symbol, interval, limit)
//
//		if macdValue < 0 && state.PrevMACDValue > 0 {
//			var botMessage models.BotMessage
//
//			botMessage = models.BotMessage{
//				ChatId: int(chatID),
//				Text:   "Значение MACD опустилось на красную отметку 🔴\n" + "Текущее значение: " + strconv.FormatFloat(macdValue, 'f', -1, 64),
//			}
//			logger.Println("Red Flag")
//
//			buf, err := json.Marshal(botMessage)
//			if err != nil {
//				logger.Println("Ошибка при маршалинге сообщения:", err)
//				continue
//			}
//			_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
//			if err != nil {
//				logger.Println("Ошибка при отправке сообщения:", err)
//			}
//			logger.Println("Red Flag notified successfully")
//		}
//
//		setPrevMACDValue(chatID, macdValue)
//		logger.Println(`Previous MACD value set - value ` + fmt.Sprint(macdValue))
//
//		time.Sleep(time.Minute * 10)
//	}
//}

// GetMACDLoopGreen logging done
func GetMACDLoopGreen(botUrl string, chatID int64, uid string) {
	state := getUserState(chatID)

	for state.IsRunning {
		logger, file := logging.CustomLog(`chatId=`+fmt.Sprint(chatID), uid)
		macdValue := GetMACD(client, symbol, interval, limit)

		if macdValue > 0 && state.PrevMACDValue <= 0 {
			var botMessage models.BotMessage

			botMessage = models.BotMessage{
				ChatId: int(chatID),
				Text:   "Значение MACD поднялось на зеленую отметку 🟢\n" + "Текущее значение: " + strconv.FormatFloat(macdValue, 'f', -1, 64),
			}

			buf, err := json.Marshal(botMessage)
			if err != nil {
				logger.Println("Ошибка при маршалинге сообщения:", err)
				continue
			}
			_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
			if err != nil {
				logger.Println("Ошибка при отправке сообщения:", err)
			}
			logger.Println("Green Flag notified successfully")
		}

		setPrevMACDValue(chatID, macdValue)
		logger.Println(`Previous MACD value set - value ` + fmt.Sprint(macdValue))

		time.Sleep(time.Minute * 10)
		file.Close()
	}
}
