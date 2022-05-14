package consts

import (
	"fmt"
	"path/filepath"
)

const PidFilePath = "%s/pid"

func GetPid(home string) string {
	return filepath.FromSlash(fmt.Sprintf(PidFilePath, home))
}

const StartupWaitTime = 2
