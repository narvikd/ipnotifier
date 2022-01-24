package main

import (
	"ipnotifier/iputils"
	"ipnotifier/telegram"
	"log"
	"os"
)

func init() {
	var envVars = []string{"token", "chatid"}
	for _, v := range envVars {
		if len(os.Getenv(v)) <= 1 {
			log.Fatalf("environment variable \"%s\" is not set.\n", v)
		}
	}
}

func main() {
	errSendIP := sendIP()
	if errSendIP != nil {
		log.Fatalln(errSendIP)
	}
}

func sendIP() error {
	ip, errIP := iputils.GetPublicIP()
	if errIP != nil {
		return errIP
	}

	m := telegram.NewClientReqModel(ip, os.Getenv("token"), os.Getenv("chatid"))
	errSend := m.Send()
	if errSend != nil {
		return errSend
	}
	return nil
}
