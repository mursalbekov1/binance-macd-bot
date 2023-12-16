package binance

import (
	"binance_tg/models"
	"bytes"
	"encoding/json"
	"github.com/adshao/go-binance/v2"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type UserState struct {
	IsRunning     bool
	PrevMACDValue float64
	IsFirstRun    bool
	IsAuthorized  bool
}

var (
	apiKey     = ""
	secretKey  = ""
	client     = binance.NewClient(apiKey, secretKey)
	symbol     = "BTCUSDT"
	password   = "0214234"
	interval   = "1s"
	limit      = 100
	userStates = make(map[int64]*UserState) // Хранение состояния для каждого пользователя
	mu         sync.Mutex
)

// Функция для получения состояния пользователя
func getUserState(chatID int64) *UserState {
	mu.Lock()
	defer mu.Unlock()

	state, ok := userStates[chatID]
	if !ok {
		// Создание нового состояния для нового пользователя
		state = &UserState{IsFirstRun: true}
		userStates[chatID] = state
	}

	return state
}

// Остальной код остается неизменным
// ...

func checkAuthorization(chatID int64) bool {
	state := getUserState(chatID)
	return state.IsAuthorized
}

func setAuthorization(chatID int64, authorized bool) {
	state := getUserState(chatID)
	state.IsAuthorized = authorized
}

func setRunning(chatID int64, running bool) {
	state := getUserState(chatID)
	state.IsRunning = running
}

func setPrevMACDValue(chatID int64, value float64) {
	state := getUserState(chatID)
	state.PrevMACDValue = value
}

func setFirstRun(chatID int64, isFirstRun bool) {
	state := getUserState(chatID)
	state.IsFirstRun = isFirstRun
}

func Respond(botUrl string, update models.Update) error {
	var botMessage models.BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId

	state := getUserState(int64(botMessage.ChatId))

	switch update.Message.Text {
	case "/start":
		if state.IsAuthorized {
			botMessage.Text = "Привет! 🌟 Добро пожаловать в MACD Notifier Bot! 📈\n\nЭтот бот предоставит вам уведомления о изменениях в значении MACD на бирже Binance.\n\n🚀 Команды: \n- /launch - Запускает бот и уведомляет вас, когда значение MACD переходит с красной зоны на зеленую и обратно.\n /red - Уведомляет вас, когда значение MACD переходит с зеленой зоны на красную.\n /green - Уведомляет вас, когда значение MACD переходит с красной зоны на зеленую.\n /stop - Останавливает текущее действие."
		} else {
			botMessage.Text = "Привет! 🌟 Добро пожаловать в MACD Notifier Bot! 📈\n\nЭтот бот предоставит вам уведомления о изменениях в значении MACD на бирже Binance. 🚀\n\nВведите пароль, чтобы получить доступ к боту: 🔐"
		}
	case "/launch":
		if !checkAuthorization(int64(botMessage.ChatId)) {
			botMessage.Text = "Сначала введите пароль, чтобы получить доступ к боту: 🔐"
			break
		}
		if state.IsRunning {
			botMessage.Text = "MACD Notifier уже запущен."
		} else {
			botMessage.Text = "MACD Notifier запущен! 📈\n\nТеперь я буду уведомлять вас о изменениях в значении MACD. 🚀\n\n"
			setRunning(int64(botMessage.ChatId), true)

			if state.IsFirstRun {
				currentMACD := GetMACD(client, symbol, interval, limit)
				if currentMACD > 0 {
					botMessage.Text += "Сейчас значение MACD на зеленой отметке 🟢\n" + "Текущее значение:" + strconv.FormatFloat(currentMACD, 'f', -1, 64)
				} else {
					botMessage.Text += "Сейчас значение MACD на красной отметке 🔴\n" + "Текущее значение:" + strconv.FormatFloat(currentMACD, 'f', -1, 64)
				}
				setFirstRun(int64(botMessage.ChatId), false)
			}

			go GetMACDLoop(botUrl, int64(botMessage.ChatId))
		}
	case "/red":
		if !checkAuthorization(int64(botMessage.ChatId)) {
			botMessage.Text = "Сначала введите пароль, чтобы получить доступ к боту: 🔐"
			break
		}

		if state.IsRunning {
			botMessage.Text = "MACD Notifier уже запущен."
		} else {
			botMessage.Text = "MACD Notifier запущен! 📈\n\nТеперь я буду уведомлять вас о изменениях в значении MACD когда он переходит с зеленой зоны на красную. 🚀\n\n"
			setRunning(int64(botMessage.ChatId), true)

			if state.IsFirstRun {
				currentMACD := GetMACD(client, symbol, interval, limit)
				if currentMACD < 0 {
					botMessage.Text += "Сейчас значение MACD на красной отметке 🔴\n" + "Текущее значение: " + strconv.FormatFloat(currentMACD, 'f', -1, 64)
				} else {
					botMessage.Text += "Сейчас значение MACD на зеленой отметке 🟢\n" + "Текущее значение: " + strconv.FormatFloat(currentMACD, 'f', -1, 64)
				}
				setFirstRun(int64(botMessage.ChatId), false)
			}

			go GetMACDLoopRed(botUrl, int64(botMessage.ChatId))
		}
	case "/green":
		if !checkAuthorization(int64(botMessage.ChatId)) {
			botMessage.Text = "Сначала введите пароль, чтобы получить доступ к боту: 🔐"
			break
		}

		if state.IsRunning {
			botMessage.Text = "MACD Notifier уже запущен."
		} else {
			botMessage.Text = "MACD Notifier запущен! 📈\n\nТеперь я буду уведомлять вас о изменениях в значении MACD когда он переходит с красной зоны на зеленую. 🚀\n\n"
			setRunning(int64(botMessage.ChatId), true)

			if state.IsFirstRun {
				currentMACD := GetMACD(client, symbol, interval, limit)
				if currentMACD > 0 {
					botMessage.Text += "Сейчас значение MACD на зеленой отметке 🟢\n" + "Текущее значение: " + strconv.FormatFloat(currentMACD, 'f', -1, 64)
				} else {
					botMessage.Text += "Сейчас значение MACD на красной отметке 🔴\n" + "Текущее значение: " + strconv.FormatFloat(currentMACD, 'f', -1, 64)
				}
				setFirstRun(int64(botMessage.ChatId), false)
			}

			go GetMACDLoopGreen(botUrl, int64(botMessage.ChatId))
		}
	case "/stop":
		if state.IsRunning {
			setRunning(int64(botMessage.ChatId), false)
			botMessage.Text = "MACD Notifier остановлен."
			setFirstRun(int64(botMessage.ChatId), true)
		} else if !checkAuthorization(int64(botMessage.ChatId)) {
			botMessage.Text = "Сначала введите пароль, чтобы получить доступ к боту: 🔐"
			break
		} else {
			botMessage.Text = "MACD Notifier не запущен чтобы остановить."
		}
	default:
		if !checkAuthorization(int64(botMessage.ChatId)) && update.Message.Text == password {
			setAuthorization(int64(botMessage.ChatId), true)
			go func() {
				time.Sleep(5 * time.Minute)
				setAuthorization(int64(botMessage.ChatId), false)
				log.Println("Авторизация сброшена.")
			}()
			botMessage.Text = "✅ Пароль принят! \n\n 🚀 Доступные команды: \n /launch - Запускает бот и уведомляет вас, когда значение MACD переходит с красной зоны на зеленую и обратно.\n /red - Уведомляет вас, когда значение MACD переходит с зеленой зоны на красную.\n /green - Уведомляет вас, когда значение MACD переходит с красной зоны на зеленую.\n /stop - Останавливает текущее действие."
		} else if !checkAuthorization(int64(botMessage.ChatId)) {
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
	state := getUserState(chatID)

	for state.IsRunning {
		macdValue := GetMACD(client, symbol, interval, limit)

		if (macdValue > 0 && state.PrevMACDValue <= 0) || (macdValue <= 0 && state.PrevMACDValue > 0) {
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

		setPrevMACDValue(chatID, macdValue)

		time.Sleep(time.Second / 2)
	}
}

func GetMACDLoopRed(botUrl string, chatID int64) {
	state := getUserState(chatID)

	for state.IsRunning {
		macdValue := GetMACD(client, symbol, interval, limit)

		if macdValue < 0 && state.PrevMACDValue > 0 {
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

		setPrevMACDValue(chatID, macdValue)

		time.Sleep(time.Second / 2)
	}
}

func GetMACDLoopGreen(botUrl string, chatID int64) {
	state := getUserState(chatID)

	for state.IsRunning {
		macdValue := GetMACD(client, symbol, interval, limit)

		if macdValue > 0 && state.PrevMACDValue <= 0 {
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

		setPrevMACDValue(chatID, macdValue)

		time.Sleep(time.Second / 2)
	}
}
