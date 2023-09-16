package api

import (
	"binance_tg/internal/models"
	"bytes"
	"encoding/json"
	"github.com/adshao/go-binance/v2"
	"net/http"
	"strconv"
)

var (
	apiKey    = "DQbMjZ54eTkw7pHIxYYW1UUFWNADxaETaE79C6Ad0VH69QImdQqVLE9rpJ6ZFc28"
	secretKey = "Mdiyex2E6kiQf2jmOSjrXKsTbqTb4SdURt3LqumbqZi3OdNSUhP3I0NTi8dHdBDG"
	client    = binance.NewClient(apiKey, secretKey)
	symbol    = "BTCUSDT"
	interval  = "1s"
	limit     = 100
)

func Respond(botUrl string, update models.Update) error {
	var botMessage models.BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	botMessage.Text = strconv.FormatFloat(GetMACD(client, symbol, interval, limit), 'f', -1, 64)
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
