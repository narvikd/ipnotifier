package main

import (
	"log"
	"os"
	"reflect"
	"strings"
)

type model struct {
	token      string
	chat       string
	machine    string
	machineSet bool
}

func newModel() *model {
	m := model{
		token:      os.Getenv("token"),
		chat:       os.Getenv("chat"),
		machine:    os.Getenv("machine"),
		machineSet: true,
	}

	// check if the environment variables are set, if they aren't it exists the application with an error.
	v := reflect.ValueOf(&m).Elem() // Gets the value of the pointer model
	for i := 0; i < v.NumField(); i++ {
		itemName := strings.ToLower(v.Type().Field(i).Name)
		varValue := v.Field(i).String()
		if varValue == "" && !strings.HasPrefix(itemName, "machine") {
			log.Fatalf("environment variable \"%s\" is not set.\n", itemName)
		}
	}

	if m.machine == "" {
		m.machineSet = false
	}

	return &m
}
