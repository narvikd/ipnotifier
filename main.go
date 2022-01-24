package main

import (
	"fmt"
	"ipnotifier/fileio"
	"ipnotifier/iputils"
	"ipnotifier/telegram"
	"log"
	"os"
	"reflect"
	"strings"
)

type Model struct {
	MachineID string
	Token     string
	ChatID    string
	IPFile    string
}

func main() {
	m := Model{
		MachineID: os.Getenv("machineid"),
		Token:     os.Getenv("token"),
		ChatID:    os.Getenv("chatid"),
		IPFile:    "ip.txt",
	}
	checkEnv(&m)
	errSendIP := sendIP(&m)
	if errSendIP != nil {
		log.Fatalln(errSendIP)
	}
}

func checkEnv(m *Model) {
	v := reflect.ValueOf(m).Elem() // Gets the value of the pointer model
	for i := 0; i < v.NumField(); i++ {
		itemName := v.Type().Field(i).Name
		varValue := v.Field(i).String()
		if varValue == "" {
			log.Fatalf("environment variable \"%s\" is not set.\n", strings.ToLower(itemName))
		}
	}
}

func sendIP(m *Model) error {
	oldIP, errOldIP := fileio.ReadIP(m.IPFile)
	if errOldIP != nil {
		return errOldIP
	}

	newIP, errIP := iputils.GetPublicIP()
	if errIP != nil {
		return errIP
	}

	if oldIP != newIP {
		msg := fmt.Sprintf("Machine: %s.\nNew ip: %s", m.MachineID, newIP)

		tele := telegram.NewClientReqModel(msg, m.Token, m.ChatID)
		errSend := tele.Send()
		if errSend != nil {
			return errSend
		}

		errWriteIP := fileio.WriteIP(newIP, m.IPFile)
		if errWriteIP != nil {
			return errWriteIP
		}
	}

	return nil
}
