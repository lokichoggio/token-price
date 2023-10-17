package main

import (
	"log"

	"token-price/cmd/app"
)

func main() {
	command := app.NewScanCommand()
	if err := command.Execute(); err != nil {
		log.Fatalf("cmd Execute error: %s", err)
	}
}
