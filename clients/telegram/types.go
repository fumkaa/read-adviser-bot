package telegram

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type From struct {
	UserName string `json:"username"`
}

type Chat struct {
	ChatId int `json:"id"`
}
