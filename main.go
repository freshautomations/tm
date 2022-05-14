package main

import (
	"log"
	"tm/tm/v2/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
