package main

import (
	"log"

	"github.com/MaxGolubev19/GoCalculator/internal/application"
)

func main() {
	app := application.New()
	err := app.Run()
	if err != nil {
		log.Println(err)
	}
}
