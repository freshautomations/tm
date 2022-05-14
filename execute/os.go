package execute

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
	"tm/m/v2/consts"
	"tm/m/v2/ux"
)

func GetPid(home string) *int {
	pidFile := consts.GetPid(home)
	if _, err := os.Stat(pidFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		ux.Fatal("could not check %s", pidFile)
	}
	if bytes, err := ioutil.ReadFile(pidFile); err == nil {
		var pid int
		pid, err = strconv.Atoi(strings.Split(string(bytes), "\n")[0])
		if err != nil {
			ux.Debug("invalid data in file %s", pidFile)
			_ = os.Remove(pidFile)
			return nil
		}
		var process *os.Process
		process, err = os.FindProcess(pid)
		if err != nil {
			ux.Debug("could not query process ID %d for %s", pid, pidFile)
			_ = os.Remove(pidFile)
			return nil
		}
		err = process.Signal(syscall.Signal(0))
		if err != nil {
			if errors.Is(err, syscall.EPERM) {
				ux.Debug("user does not own process %d", pid)
				_ = os.Remove(pidFile)
				return nil
			}
			if errors.Is(err, os.ErrProcessDone) {
				ux.Debug("process %d is already done", pid)
				_ = os.Remove(pidFile)
				return nil
			}
			ux.Debug("checking if process %d is running failed: %s", pid, err)
			return nil
		}
		return &pid
	}
	ux.Fatal("could not read %s", pidFile)
	return nil
}
