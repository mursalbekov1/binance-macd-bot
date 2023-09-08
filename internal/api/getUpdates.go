package api

import (
	"binance_tg/internal/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetUpdates(botUrl string, offset int) ([]models.Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
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
