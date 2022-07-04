package main

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat     Chat   `json:"chat"`
	Text     string `json:"text"`
	FromUser User   `json:"from"`
}

type User struct {
	Id int `json:"id"`
}

type Chat struct {
	ChatId int `json:"id"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}

type BotMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}
