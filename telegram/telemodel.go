package telegram

type TeleReqModel struct {
	ChatId              string `json:"chat_id"`
	Text                string `json:"text"`
	DisableNotification bool   `json:"disable_notification"`
}

type TeleResModel struct {
	Ok bool `json:"ok"`
}

func NewTeleReqModel(msg string, chatID string) *TeleReqModel {
	model := TeleReqModel{
		ChatId:              chatID,
		Text:                msg,
		DisableNotification: true,
	}
	return &model
}
