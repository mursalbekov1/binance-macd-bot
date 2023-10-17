package api

import (
	"binance_tg/internal/models"
	"bytes"
	"encoding/json"
	"github.com/adshao/go-binance/v2"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	apiKey    = ""
	secretKey = " "
	client    = binance.NewClient(apiKey, secretKey)
	symbol    = "BTCUSDT"
	interval  = "1s"
	limit     = 100
	isRunning = false
)

func Respond(botUrl string, update models.Update) error {
	var botMessage models.BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	if update.Message.Text == "/start" {
		if isRunning {
			botMessage.Text = "MACD уже запущено."
		} else {
			botMessage.Text = "MACD запущено."
			isRunning = true
			go GetMACDLoop(botUrl, int64(botMessage.ChatId)) // Запускаем GetMACD в отдельной горутине
		}
	} else if update.Message.Text == "/stop" {
		if isRunning {
			botMessage.Text = "MACD остановлено."
			isRunning = false
		} else {
			botMessage.Text = "MACD уже остановлено."
		}
	} else {
		botMessage.Text = "Неизвестная команда."
	}

	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}

func GetMACDLoop(botUrl string, chatID int64) {
	for isRunning {
		macdValue := GetMACD(client, symbol, interval, limit)
		var botMessage models.BotMessage
		if macdValue > 0 {
			botMessage = models.BotMessage{
				ChatId: int(chatID),
				Text:   strconv.FormatFloat(macdValue, 'f', -1, 64) + " 🟢",
			}
		} else {
			botMessage = models.BotMessage{
				ChatId: int(chatID),
				Text:   strconv.FormatFloat(macdValue, 'f', -1, 64) + " 🔴",
			}
		}
		buf, err := json.Marshal(botMessage)
		if err != nil {
			log.Println("Ошибка при маршалинге сообщения:", err)
			continue
		}
		_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
		if err != nil {
			log.Println("Ошибка при отправке сообщения:", err)
		}

		time.Sleep(time.Second) // Подождать 1 секунду перед следующей проверкой
	}
}
