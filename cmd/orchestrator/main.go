package main

import (
	"log"

	"github.com/MaxGolubev19/GoCalculator/internal/orchestrator"
)

func main() {
	orch := orchestrator.New()
	if err := orch.Run(); err != nil {
		log.Println(err)
	}
}
