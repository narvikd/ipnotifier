package main

import (
	"errors"
	"fmt"
	"ipnotifier/pkg/errorsutils"
	"ipnotifier/pkg/fileutils"
	"ipnotifier/pkg/iputils"
	"ipnotifier/telegram"
	"log"
	"time"
)

func main() {
	m := newModel()
	log.Println("IP Notifier starting...")

	for {
		log.Printf("Checking if the machine: %s has a new IP Address", m.machine)
		sent, errSendIP := sendIP(m)
		if errSendIP != nil {
			log.Println(errSendIP)
		}

		if sent {
			log.Println("Sent message with new IP Address")
		} else {
			log.Println("No new IP was found")
		}

		time.Sleep(10 * time.Minute)
	}
}

// sendIP returns true if a message has been sent, or false if it hasn't been sent
func sendIP(m *model) (bool, error) {
	const ipPath = "ip.txt"
	oldIP, errOldIP := readIP(ipPath)
	if errOldIP != nil {
		return false, errOldIP
	}

	newIP, errIP := iputils.GetPublicIP()
	if errIP != nil {
		return false, errIP
	}

	if oldIP != newIP {
		msg := fmt.Sprintf("Machine: %s.\nNew ip: %s", m.machine, newIP)
		if !m.machineSet {
			msg = fmt.Sprintf("New ip: %s", newIP)
		}

		errSend := telegram.Send(telegram.NewClientModel(msg, m.token, m.chat))
		if errSend != nil {
			return false, errSend
		}

		errWriteIP := fileutils.Write(newIP, ipPath)
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
