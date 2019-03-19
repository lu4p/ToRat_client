package client

import (
	"os"
	"path/filepath"
)

var (
	Path         = filepath.Join(os.ExpandEnv("$lib"), "MacSafe")
	PathExe      = filepath.Join(Path, "MacSafe")
	TempZip      = filepath.Join(Path, "temp.zip")
	TorPath      = filepath.Join(Path, "Helper")
	TorExe       = filepath.Join(TorPath, "MacSafeHelper")
	TorData      = filepath.Join(Path, "Helper", "conf")
	HostnamePath = filepath.Join(Path, "token")
)
