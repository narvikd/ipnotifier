package main

import (
	"fmt"
	"ipnotifier/fileio"
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
	oldIP, errOldIP := fileio.ReadIP()
	if errOldIP != nil {
		return errOldIP
	}

	newIP, errIP := iputils.GetPublicIP()
	if errIP != nil {
		return errIP
	}

	if oldIP != newIP {
		msg := fmt.Sprintf("old ip: %s. new ip: %s", oldIP, newIP)
		if oldIP == "" {
			msg = fmt.Sprintf("new ip: %s", newIP)
		}

		m := telegram.NewClientReqModel(msg, os.Getenv("token"), os.Getenv("chatid"))
		errSend := m.Send()
		if errSend != nil {
			return errSend
		}

		errWriteIP := fileio.WriteIP(newIP)
		if errWriteIP != nil {
			return errWriteIP
		}
	}

	return nil
}
