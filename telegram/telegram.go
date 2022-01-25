package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"ipnotifier/pkg/errorsutils"
	"ipnotifier/pkg/httpclient"
	"net/http"
	"time"
)

type ClientModel struct {
	Client         *http.Client
	Endpoint       string
	Token          string
	ChatID         string
	TimeoutSeconds time.Duration
	Message        string
	RequestModel   *RequestModel
}

func NewClientModel(message string, token string, chatId string) *ClientModel {
	m := ClientModel{
		Client:       httpclient.MakeDefaultClient(),
		Endpoint:     "https://api.telegram.org/bot" + token + "/sendMessage",
		Token:        token,
		ChatID:       chatId,
		Message:      message,
		RequestModel: newRequestModel(message, chatId),
	}
	return &m
}

func Send(m *ClientModel) error {
	reqBody, errMarshal := json.Marshal(m.RequestModel)
	if errMarshal != nil {
		return errorsutils.Wrap(errMarshal, "telegram http client couldn't marshall telegram model")
	}

	req, errRes := m.Client.Post(m.Endpoint, "application/json", bytes.NewBuffer(reqBody))
	if errRes != nil {
		return errorsutils.Wrap(errRes, "telegram http client couldn't make http request")
	}
	//req.Header.Add("User-Agent", "")

	defer req.Body.Close()
	if req.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram http client responded with a non 200 code: %v", req.StatusCode)
	}

	resModel := new(ResponseModel)
	errResDecode := json.NewDecoder(req.Body).Decode(resModel)
	if errResDecode != nil {
		return errorsutils.Wrap(errResDecode, "telegram http client couldn't decode response into response model")
	}

	if !resModel.Ok {
		return errors.New("telegram http client received a non-ok response from telegram's api")
	}
	return nil
}
