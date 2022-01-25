package main

import (
	"errors"
	"fmt"
	"ipnotifier/pkg/convert"
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

type model struct {
	machineID   string
	token       string
	chatID      string
	ipFile      string
	refreshmins string
	refreshTime time.Duration
}

func main() {
	m := model{
		machineID:   os.Getenv("machineid"),
		token:       os.Getenv("token"),
		chatID:      os.Getenv("chatid"),
		ipFile:      "ip.txt",
		refreshmins: os.Getenv("refreshmins"),
	}
	checkEnv(&m)
	log.Println("IP Notifier starting...")

	for {
		log.Printf("Checking if the machine: %s has a new IP Address", m.machineID)
		sent, errSendIP := sendIP(&m)
		if errSendIP != nil {
			log.Println(errSendIP)
		}

		if sent {
			log.Println("Sent message with new IP Address")
		} else {
			log.Println("No new IP was found")
		}

		time.Sleep(m.refreshTime)
	}
}

// checkEnv checks if the environment variables are set, if they aren't it exists the application with an error.
//
// "refreshTime" is not verified, but "refreshStr" is. If the conversion from str to time.Duration is not successful it also exits the app.
func checkEnv(m *model) {
	v := reflect.ValueOf(m).Elem() // Gets the value of the pointer model
	for i := 0; i < v.NumField(); i++ {
		itemName := v.Type().Field(i).Name
		varValue := v.Field(i).String()
		if varValue == "" && itemName != "refreshTime" {
			log.Fatalf("environment variable \"%s\" is not set.\n", strings.ToLower(itemName))
		}
	}

	// Sets m.refreshTime
	duration, errConvert := convert.StrToDuration(m.refreshmins)
	if errConvert != nil {
		log.Fatalln("environment variable \"refreshmins\" couldn't be converted to a time, please verify if it's set correctly.")
	}
	m.refreshTime = duration * time.Minute
}

// sendIP returns true if a message has been sent, or false if it hasn't been sent
func sendIP(m *model) (bool, error) {
	oldIP, errOldIP := readIP(m.ipFile)
	if errOldIP != nil {
		return false, errOldIP
	}

	newIP, errIP := iputils.GetPublicIP()
	if errIP != nil {
		return false, errIP
	}

	if oldIP != newIP {
		msg := fmt.Sprintf("Machine: %s.\nNew ip: %s", m.machineID, newIP)

		tele := telegram.NewClientReqModel(msg, m.token, m.chatID)
		errSend := tele.Send()
		if errSend != nil {
			return false, errSend
		}

		errWriteIP := fileutils.Write(newIP, m.ipFile)
		if errWriteIP != nil {
			return false, errWriteIP
		}

		return true, nil
	}

	return false, nil
}

// readIP reads the ip file contents from the path and returns an ip as a string.
//
// If the file contents are empty, it needs to return an empty string, not an error,
// since it has to verify if the old ip equals the old ip (if there's an error, it cannot verify it).
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
