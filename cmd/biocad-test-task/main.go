package main

import (
	"log"
	"test_task/internal/pkg/app"
)

func main() {
	log.Print("start")

	App, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	err = App.Run()
	if err != nil {
		log.Fatal(err)
	}
}
