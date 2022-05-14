package utils

import (
	"io/fs"
	"os"
	"tm/tm/v2/ux"
)

func WriteFile(path string, content string) {
	if err := os.WriteFile(path, []byte(content), fs.ModePerm); err != nil {
		ux.Fatal(err.Error())
	}
}
