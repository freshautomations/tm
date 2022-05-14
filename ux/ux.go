package ux

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

func FatalRaw(format string, a ...any) {
	log.New(os.Stderr, "", 0).Fatalf(format, a...)
}

func Fatal(format string, a ...any) {
	FatalRaw(fmt.Sprintf("Error: %s.", format), a...)
}

func WarnRaw(format string, a ...any) {
	if !viper.GetBool("quiet") {
		log.New(os.Stderr, "", 0).Printf(format, a...)
	}
}

func Warn(format string, a ...any) {
	WarnRaw(fmt.Sprintf("Warning: %s.", format), a...)
}

func Debug(format string, a ...any) {
	if viper.GetBool("debug") {
		log.New(os.Stderr, "DEBUG ", 0).Printf(format, a...)
	}
}

func Info(format string, a ...any) {
	if !viper.GetBool("quiet") {
		log.New(os.Stdout, "", 0).Printf(format, a...)
	}
}
