// +build notor

package client

import "log"

func installTor() error {
	log.Println("not installing tor")
	return nil
}
