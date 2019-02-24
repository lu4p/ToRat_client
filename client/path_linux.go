package client

import (
	"os"
	"path/filepath"
)

var (
	Path         = filepath.Join(os.ExpandEnv("$lib"), "libssh")
	PathExe      = filepath.Join(Path, "libssh")
	TempZip      = filepath.Join(Path, "temp.zip")
	TorPath      = filepath.Join(Path, "")
	TorExe       = filepath.Join(TorPath, "")
	TorData      = filepath.Join(Path, "etc", "tor")
	HostnamePath = filepath.Join(Path, "token")
)
