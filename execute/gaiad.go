package execute

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
	"tm/tm/v2/consts"
	"tm/tm/v2/utils"
	"tm/tm/v2/ux"
)

func execute(binary string, arg ...string) (string, error) {
	return executeWithStdIn("", binary, arg...)
}

func executeWithStdIn(stdin string, binary string, arg ...string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(binary, arg...)
	cmd.Stdin = strings.NewReader(stdin)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	ux.Debug("%s %s\n", binary, strings.Join(arg, " "))
	err := cmd.Run()
	if err != nil {
		return stderr.String(), err
	}
	if len(stderr.Bytes()) != 0 {
		return stderr.String(), nil
	}
	return stdout.String(), nil
}

func startupChecker(checker chan int, process *os.Process) {
	_, _ = process.Wait()
	checker <- 1
}

func executeStart(binary string, arg ...string) (int, error) {
	cmd := exec.Command(binary, arg...)
	ux.Debug("%s %s\n", binary, strings.Join(arg, " "))
	err := cmd.Start()
	if err != nil {
		return 0, err
	}
	checker := make(chan int, 1)
	go startupChecker(checker, cmd.Process)
	for i := 0; i < consts.StartupWaitTime; i++ {
		time.Sleep(time.Second)
		if len(checker) > 0 {
			return 0, fmt.Errorf("PID %d stopped", cmd.Process.Pid)
		}
	}
	return cmd.Process.Pid, nil
}

func fatal(out string, err error) {
	if err != nil {
		ErrorString := strings.Split(out, "\n")[0]
		ux.FatalRaw(ErrorString)
	}
}

func debug(out string, err error) {
	if err != nil {
		ErrorString := strings.Split(out, "\n")[0]
		ux.Debug("%s", ErrorString)
	}
}

/*
func warn(out string, err error) {
	if err != nil {
		ErrorString := strings.Split(out, "\n")[0]
		ux.WarnRaw("%s", ErrorString)
	}
}
*/

func Init(fullNodename string, binary string, home string) {
	fullNodenameSplit := strings.Split(fullNodename, ".")
	chainName := fullNodenameSplit[0]
	nodeName := fullNodenameSplit[1]
	args := []string{"init", nodeName, "--chain-id", chainName, "--home", home}
	out, err := execute(binary, args...)
	switch {
	case err != nil:
		debug(out, err)
	case utils.GetConfigEntryContentString(out, "json", "moniker") != nodeName:
		ux.Warn("could not initialize %s", fullNodenameSplit)
	default:
		ux.Debug("successful init: %s", fullNodename)
	}
}

func KeysAdd(binary string, home string, name string, hdpath string, mnemonics string) {
	args := []string{"keys", "add", name, "--keyring-backend", "test", "--keyring-dir", home, "--output", "json"}
	if hdpath != "" {
		args = append(args, "--hd-path", hdpath)
	}
	if mnemonics != "" {
		args = append(args, "--recover")
	}

	output, err := executeWithStdIn(mnemonics, binary, args...)
	if err != nil {
		ux.Debug("did not add key %s to %s", name, home)
		return
	}
	utils.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/config/mnemonics/%s.json", home, name)), output)
	ux.Debug("successful key add %s to chain %s", name, home)
}

func AddGenesisAccount(binary string, home string, name string, value string) {
	args := []string{"add-genesis-account", name, value, "--keyring-backend", "test", "--home", home, "--output", "json"}

	out, err := execute(binary, args...)
	debug(out, err)
	if err == nil {
		ux.Debug("successful genesis account addition %s to %s", name, home)
	}
}

func AddGentx(binary string, home string, chainID string, name string, value string) {
	args := []string{"gentx", name, value, "--keyring-backend", "test", "--home", home, "--moniker", name, "--chain-id", chainID, "--output", "json"}

	out, err := execute(binary, args...)
	if err != nil {
		ErrorString := strings.Split(out, "\n")[0]
		if strings.HasPrefix(ErrorString, "Error: failed to write signed gen tx: open") &&
			strings.HasSuffix(ErrorString, ": file exists") {
			ux.Debug("gentx already exists for %s in %s", name, home)
		} else {
			fatal(out, err)
		}
	}
	if err == nil {
		ux.Debug("successful gentx addition %s to %s", name, home)
	}
}

func CollectGentxs(binary string, home string) {
	args := []string{"collect-gentxs", "--home", home}

	_, err := execute(binary, args...)
	if err != nil {
		ux.Fatal("could not collect gentxs %s", home)
		return
	}
	ux.Debug("successful gentx collection for %s", home)
}

func ValidateGenesis(binary string, home string) {
	args := []string{"validate-genesis", "--home", home}

	_, err := execute(binary, args...)
	if err != nil {
		ux.Fatal("could not validate genesis %s", home)
		return
	}
	ux.Debug("successful validation of genesis for %s", home)
}

func ShowNodeID(binary string, home string) string {
	args := []string{"tendermint", "show-node-id", "--home", home}

	output, err := execute(binary, args...)
	fatal(output, err)
	output = strings.Split(output, "\n")[0]
	if output == "" {
		ux.Fatal("could not retrieve node ID from %s", home)
	}
	return output
}

func Start(binary string, home string) (int, error) {
	args := []string{"start", "--home", home}

	pid, err := executeStart(binary, args...)
	if err != nil {
		return 0, err
	}
	pidString := strconv.Itoa(pid)
	err = ioutil.WriteFile(consts.GetPid(home), []byte(pidString), fs.ModePerm)
	ux.Debug("process %s started and written to %s", pidString, consts.GetPid(home))
	return pid, err
}

func Stop(home string) error {
	pid := GetPid(home)
	process, err := os.FindProcess(*pid)
	if err != nil {
		_ = GetPid(home)
		return nil
	}
	err = process.Signal(syscall.SIGINT)
	if err != nil {
		return fmt.Errorf("could not SIGINT PID %d", pid)
	}
	_ = GetPid(home)
	return nil
}
