package telegram

// RequestModel is the model that the Telegram API accepts to send a message via the bot.
type RequestModel struct {
	ChatId              string `json:"chat_id"`
	Text                string `json:"text"`
	DisableNotification bool   `json:"disable_notification"`
}

// ResponseModel is the model that the Telegram API returns when RequestModel is received by it.
type ResponseModel struct {
	Ok bool `json:"ok"`
}

// newRequestModel returns a new pointer of RequestModel, with DisableNotification set to true .
func newRequestModel(msg string, chatID string) *RequestModel {
	model := RequestModel{
		ChatId:              chatID,
		Text:                msg,
		DisableNotification: true,
	}
	return &model
}
