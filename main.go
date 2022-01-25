package main

import (
	"errors"
	"fmt"
	"ipnotifier/pkg/errorsutils"
	"ipnotifier/pkg/fileutils"
	"ipnotifier/pkg/iputils"
	"ipnotifier/telegram"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
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
	log.Println("IP Notifier starting...")

	for {
		log.Printf("Checking if the machine: %s has a new IP Address", m.MachineID)
		sent, errSendIP := sendIP(&m)
		if errSendIP != nil {
			log.Println(errSendIP)
		}

		if sent {
			log.Println("Sent message with new IP Address")
		} else {
			log.Println("No new IP was found")
		}

		time.Sleep(15 * time.Minute)
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

// sendIP returns true if a message has been sent, or false if it hasn't been sent
func sendIP(m *Model) (bool, error) {
	oldIP, errOldIP := readIP(m.IPFile)
	if errOldIP != nil {
		return false, errOldIP
	}

	newIP, errIP := iputils.GetPublicIP()
	if errIP != nil {
		return false, errIP
	}

	if oldIP != newIP {
		msg := fmt.Sprintf("Machine: %s.\nNew ip: %s", m.MachineID, newIP)

		tele := telegram.NewClientReqModel(msg, m.Token, m.ChatID)
		errSend := tele.Send()
		if errSend != nil {
			return false, errSend
		}

		errWriteIP := fileutils.Write(newIP, m.IPFile)
		if errWriteIP != nil {
			return false, errWriteIP
		}

		return true, nil
	}

	return false, nil
}

func readIP(path string) (string, error) {
	contents, err := fileutils.Read(path)
	if err != nil {
		return "", errorsutils.Wrap(err, "ip file")
	}

	if len(contents) <= 0 {
		return "", nil
	}

	ip := contents[0]
	if ip != "" && !iputils.IsIPValid(ip) {
		return "", errors.New("ip file has bogus content inside")
	}
	return ip, nil
}
