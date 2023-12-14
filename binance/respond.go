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
	prevMACDValue float64
	isFirstRun    = true
	password      = "0214234"
	isAuthorized  = false
)

func checkAuthorization() bool {
	return isAuthorized
}

func Respond(botUrl string, update models.Update) error {
	var botMessage models.BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId

	switch update.Message.Text {
	case "/start":
		if isAuthorized {
			botMessage.Text = "Привет! Этот бот предоставит уведомления о изменениях в значении MACD. Для активации введите команду /launch."
		} else {
			botMessage.Text = "Привет! Этот бот предоставит уведомления о изменениях в значении MACD. Введите пароль, чтобы получить доступ к боту:"
		}
	case "/launch":
		if !checkAuthorization() {
			botMessage.Text = "🔐 Для использования бота сначала введите пароль."
			break
		}

		if isRunning {
			botMessage.Text = "MACD уже запущено."
		} else {
			botMessage.Text = "MACD запущено."
			isRunning = true

			if isFirstRun {
				currentMACD := GetMACD(client, symbol, interval, limit)
				if currentMACD > 0 {
					botMessage.Text += " Сейчас значение MACD на зеленой отметке 🟢 " + strconv.FormatFloat(currentMACD, 'f', -1, 64)
				} else {
					botMessage.Text += " Сейчас значение MACD на красной отметке 🔴 " + strconv.FormatFloat(currentMACD, 'f', -1, 64)
				}
				isFirstRun = false
			}

			go GetMACDLoop(botUrl, int64(botMessage.ChatId))
		}
	case "/stop":
		if !checkAuthorization() {
			botMessage.Text = "🔐 Для использования бота сначала введите пароль."
			break
		}

		if isRunning {
			isRunning = false
			botMessage.Text = "MACD остановлено."
			isFirstRun = false
		} else {
			botMessage.Text = "MACD уже остановлено."
		}
	default:
		if !checkAuthorization() && update.Message.Text == password {
			isAuthorized = true
			botMessage.Text = "✨ Пароль принят! Теперь у вас есть доступ к функциям бота. Для старта используйте команду /launch ✈️."
		} else if !checkAuthorization() {
			botMessage.Text = "🔒 Неверный пароль. Введите правильный пароль, чтобы открыть доступ к функциям бота."
		} else {
			botMessage.Text = "Неизвестная команда."
		}
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

		if (macdValue > 0 && prevMACDValue <= 0) || (macdValue <= 0 && prevMACDValue > 0) {
			var botMessage models.BotMessage
			if macdValue > 0 {
				botMessage = models.BotMessage{
					ChatId: int(chatID),
					Text:   "Значение MACD поднялось на зеленый уровень 🟢 " + strconv.FormatFloat(macdValue, 'f', -1, 64),
				}
			} else {
				botMessage = models.BotMessage{
					ChatId: int(chatID),
					Text:   "Значение MACD опустилось на красный уровень 🔴 " + strconv.FormatFloat(macdValue, 'f', -1, 64),
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

		time.Sleep(time.Second / 2)
	}
}
