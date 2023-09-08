package api

import (
	"binance_tg/internal/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetUpdates(botUrl string) ([]models.Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse models.RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, err
}
