package consts

import (
	"tm/tm/v2/utils"
)

const PidFilePath = "%s/pid"
const LogFilePath = "%s/log"
const MnemonicsDirPath = "%s/config/mnemonics"
const MnemonicsPath = "%s/config/mnemonics/%s.json"

func GetPid(home string) string {
	return utils.GetSlashPath(PidFilePath, home)
}

func GetLog(home string) string {
	return utils.GetSlashPath(LogFilePath, home)
}

func GetMnemonicsDir(home string) string {
	return utils.GetSlashPath(MnemonicsDirPath, home)
}

func GetMnemonics(home string, shortNodeName string) string {
	return utils.GetSlashPath(MnemonicsPath, home, shortNodeName)
}

const StartupWaitTime = 2
