package models

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	chat Chat   `json:"chat"`
	text string `json:"text"`
}

type Chat struct {
	ChatId int `json:"chat"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}
