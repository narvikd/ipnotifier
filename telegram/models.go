package telegram

type RequestModel struct {
	ChatId              string `json:"chat_id"`
	Text                string `json:"text"`
	DisableNotification bool   `json:"disable_notification"`
}

type ResponseModel struct {
	Ok bool `json:"ok"`
}

func newRequestModel(msg string, chatID string) *RequestModel {
	model := RequestModel{
		ChatId:              chatID,
		Text:                msg,
		DisableNotification: true,
	}
	return &model
}
