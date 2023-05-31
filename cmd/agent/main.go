package main

import (
	"github.com/h2p2f/practicum-metrics/internal/client/app"
	"github.com/h2p2f/practicum-metrics/internal/logger"
	"log"
)

// function to monitor metrics

func main() {
	//init logger
	if err := logger.InitLogger("info"); err != nil {
		log.Fatal(err)
	}
	app.Run(logger.Log)
}
