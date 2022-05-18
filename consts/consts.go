package consts

import (
	"tm/tm/v2/utils"
)

const PidFilePath = "%s/pid"
const LogFilePath = "%s/log"

func GetPid(home string) string {
	return utils.GetSlashPath(PidFilePath, home)
}

func GetLog(home string) string {
	return utils.GetSlashPath(LogFilePath, home)
}

const StartupWaitTime = 2
