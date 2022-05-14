package utils

import (
	"mvdan.cc/sh/v3/shell"
	"os/exec"
	"path/filepath"
	"tm/tm/v2/ux"
)

func FindOSBinary(name string) (result string) {
	expanded, err := shell.Expand(name, nil)
	if err != nil {
		ux.Fatal("%s cannot be expanded ", name)
	}
	result, err = exec.LookPath(expanded)
	if err == nil {
		result, _ = filepath.Abs(result)
	}
	if result == "" {
		result = expanded
	}
	return
}
