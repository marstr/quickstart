package main

import (
	"log"

	"github.com/marstr/quickstart/actions"
)

func main() {
	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
