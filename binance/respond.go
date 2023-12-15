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
	password      = "0214234"
	interval      = "1s"
	limit         = 100
	isRunning     = false
	prevMACDValue float64
	isFirstRun    = true
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
			botMessage.Text = "Привет! 🌟 Добро пожаловать в MACD Notifier Bot! 📈\n\nЭтот бот предоставит вам уведомления о изменениях в значении MACD на бирже Binance.\n\n🚀 Команды: \n- /launch - Запускает бот и уведомляет вас, когда значение MACD переходит с красной зоны на зеленую и обратно.\n /red - Уведомляет вас, когда значение MACD переходит с зеленой зоны на красную.\n /green - Уведомляет вас, когда значение MACD переходит с красной зоны на зеленую.\n /stop - Останавливает текущее действие."
		} else {
			botMessage.Text = "Привет! 🌟 Добро пожаловать в MACD Notifier Bot! 📈\n\nЭтот бот предоставит вам уведомления о изменениях в значении MACD на бирже Binance. 🚀\n\nВведите пароль, чтобы получить доступ к боту: 🔐"
		}
	case "/launch":
		if !checkAuthorization() {
			botMessage.Text = "Сначала введите пароль, чтобы получить доступ к боту: 🔐"
			break
		}
		if isRunning {
			botMessage.Text = "MACD Notifier уже запущен."
		} else {
			botMessage.Text = "MACD Notifier запущен! 📈\n\nТеперь я буду уведомлять вас о изменениях в значении MACD. 🚀\n\n"
			isRunning = true

			if isFirstRun {
				currentMACD := GetMACD(client, symbol, interval, limit)
				if currentMACD > 0 {
					botMessage.Text += "Сейчас значение MACD на зеленой отметке 🟢\n" + "Текущее значение:" + strconv.FormatFloat(currentMACD, 'f', -1, 64)
				} else {
					botMessage.Text += "Сейчас значение MACD на красной отметке 🔴\n" + "Текущее значение:" + strconv.FormatFloat(currentMACD, 'f', -1, 64)
				}
				isFirstRun = false
			}

			go GetMACDLoop(botUrl, int64(botMessage.ChatId))
		}
	case "/red":
		if !checkAuthorization() {
			botMessage.Text = "Сначала введите пароль, чтобы получить доступ к боту: 🔐"
			break
		}

		if isRunning {
			botMessage.Text = "MACD Notifier уже запущен."
		} else {
			botMessage.Text = "MACD Notifier запущен! 📈\n\nТеперь я буду уведомлять вас о изменениях в значении MACD когда он переходит с зеленой зоны на красную. 🚀\n\n"
			isRunning = true

			if isFirstRun {
				currentMACD := GetMACD(client, symbol, interval, limit)
				if currentMACD < 0 {
					botMessage.Text += "Сейчас значение MACD на красной отметке 🔴\n" + "Текущее значение: " + strconv.FormatFloat(currentMACD, 'f', -1, 64)
				} else {
					botMessage.Text += "Сейчас значение MACD на зеленой отметке 🟢\n" + "Текущее значение: " + strconv.FormatFloat(currentMACD, 'f', -1, 64)
				}
				isFirstRun = false
			}

			go GetMACDLoopRed(botUrl, int64(botMessage.ChatId))
		}
	case "/green":
		if !checkAuthorization() {
			botMessage.Text = "Сначала введите пароль, чтобы получить доступ к боту: 🔐"
			break
		}

		if isRunning {
			botMessage.Text = "MACD Notifier уже запущен."
		} else {
			botMessage.Text = "MACD Notifier запущен! 📈\n\nТеперь я буду уведомлять вас о изменениях в значении MACD когда он переходит с зеленой зоны на красную. 🚀\n\n"
			isRunning = true

			if isFirstRun {
				currentMACD := GetMACD(client, symbol, interval, limit)
				if currentMACD > 0 {
					botMessage.Text += "Сейчас значение MACD на зеленой отметке 🟢\n" + "Текущее значение: " + strconv.FormatFloat(currentMACD, 'f', -1, 64)
				} else {
					botMessage.Text += "Сейчас значение MACD на красной отметке 🔴\n" + "Текущее значение: " + strconv.FormatFloat(currentMACD, 'f', -1, 64)
				}
				isFirstRun = false
			}

			go GetMACDLoopGreen(botUrl, int64(botMessage.ChatId))
		}
	case "/stop":
		if isRunning {
			isRunning = false
			botMessage.Text = "MACD Notifier остановлен."
			isFirstRun = false
		} else if !checkAuthorization() {
			botMessage.Text = "Сначала введите пароль, чтобы получить доступ к боту: 🔐"
			break
		} else {
			botMessage.Text = "MACD Notifier не запущен чтобы остановить."
		}
	default:
		if !checkAuthorization() && update.Message.Text == password {
			isAuthorized = true
			go func() {
				time.Sleep(5 * time.Minute)
				isAuthorized = false
				log.Println("Авторизация сброшена.")
			}()
			botMessage.Text = "✅ Пароль принят! \n\n 🚀 Доступные команды: \n /launch - Запускает бот и уведомляет вас, когда значение MACD переходит с красной зоны на зеленую и обратно.\n /red - Уведомляет вас, когда значение MACD переходит с зеленой зоны на красную.\n /green - Уведомляет вас, когда значение MACD переходит с красной зоны на зеленую.\n /stop - Останавливает текущее действие."
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
					Text:   "Значение MACD поднялось на зеленую отметку 🟢 \n" + "Текущее значение: " + strconv.FormatFloat(macdValue, 'f', -1, 64),
				}
			} else {
				botMessage = models.BotMessage{
					ChatId: int(chatID),
					Text:   "Значение MACD опустилось на красную отметку 🔴 \n" + "Текущее значение: " + strconv.FormatFloat(macdValue, 'f', -1, 64),
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

func GetMACDLoopRed(botUrl string, chatID int64) {
	for isRunning {
		macdValue := GetMACD(client, symbol, interval, limit)

		if macdValue < 0 && prevMACDValue > 0 {
			var botMessage models.BotMessage

			botMessage = models.BotMessage{
				ChatId: int(chatID),
				Text:   "Значение MACD опустилось на красную отметку 🔴\n" + "Текущее значение: " + strconv.FormatFloat(macdValue, 'f', -1, 64),
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

func GetMACDLoopGreen(botUrl string, chatID int64) {
	for isRunning {
		macdValue := GetMACD(client, symbol, interval, limit)

		if macdValue > 0 && prevMACDValue <= 0 {
			var botMessage models.BotMessage

			botMessage = models.BotMessage{
				ChatId: int(chatID),
				Text:   "Значение MACD поднялось на зеленую отметку 🟢\n" + "Текущее значение: " + strconv.FormatFloat(macdValue, 'f', -1, 64),
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
