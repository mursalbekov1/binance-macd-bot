package binance

import (
	"binance_tg/models"
	"bytes"
	"encoding/json"
	"github.com/adshao/go-binance/v2"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	apiKey        = ""
	secretKey     = ""
	client        = binance.NewClient(apiKey, secretKey)
	symbol        = "BTCUSDT"
	interval      = "1s"
	limit         = 100
	isRunning     = false
	lastSign      = 0
	prevMACDValue float64
	isFirstRun    = true
)

func Respond(botUrl string, update models.Update) error {
	var botMessage models.BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId

	switch update.Message.Text {
	case "/start":
		botMessage.Text = "Привет! Этот бот предоставит уведомления о изменениях в значении MACD. Для активации введи команду /launch."
	case "/launch":
		if isRunning {
			botMessage.Text = "MACD уже запущено."
		} else {
			botMessage.Text = "MACD запущено."
			isRunning = true

			if isFirstRun {
				// Отправить сообщение с текущим состоянием MACD только при первом запуске
				currentMACD := GetMACD(client, symbol, interval, limit)
				if currentMACD > 0 {
					botMessage.Text += " Сейчас значение MACD на зеленой отметке 🟢 " + strconv.FormatFloat(currentMACD, 'f', -1, 64)
				} else {
					botMessage.Text += " Сейчас значение MACD на красной отметке 🔴 " + strconv.FormatFloat(currentMACD, 'f', -1, 64)
				}
				isFirstRun = false
			}

			go GetMACDLoop(botUrl, int64(botMessage.ChatId)) // Запускаем GetMACD в отдельной горутине
		}
	case "/stop":
		if isRunning {
			isRunning = false
			botMessage.Text = "MACD остановлено."
			isFirstRun = false
		} else {
			botMessage.Text = "MACD уже остановлено."
		}
	default:
		botMessage.Text = "Неизвестная команда"
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

		// Проверяем, изменился ли знак MACD
		if (macdValue > 0 && prevMACDValue <= 0) || (macdValue <= 0 && prevMACDValue > 0) {
			var botMessage models.BotMessage
			if macdValue > 0 {
				botMessage = models.BotMessage{
					ChatId: int(chatID),
					Text:   "Значение macd поднялся на зеленый уровень 🟢 " + strconv.FormatFloat(macdValue, 'f', -1, 64),
				}
			} else {
				botMessage = models.BotMessage{
					ChatId: int(chatID),
					Text:   "Значение macd опустился на красный уровень 🔴 " + strconv.FormatFloat(macdValue, 'f', -1, 64),
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
		}

		prevMACDValue = macdValue

		time.Sleep(time.Second / 2) // Подождать 1 секунду перед следующей проверкой
	}
}
