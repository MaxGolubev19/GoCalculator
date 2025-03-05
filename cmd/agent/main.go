package main

import (
	"log"

	"github.com/MaxGolubev19/GoCalculator/internal/agent"
)

func main() {
	agent := agent.New()
	if err := agent.Run(); err != nil {
		log.Println(err)
	}
}
