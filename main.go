package main

import (
	"log"
	"tm/m/v2/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
