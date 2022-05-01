package ux

import (
	"fmt"
	"log"
	"os"
)

func Fatal(format string, a ...any) {
	log.New(os.Stderr, "", 0).Fatalf(fmt.Sprintf("Error: %s.", format), a...)
}
