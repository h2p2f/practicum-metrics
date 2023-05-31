package main

import (
	"fmt"
	"github.com/h2p2f/practicum-metrics/internal/logger"
	"github.com/h2p2f/practicum-metrics/internal/server/app"
)

func main() {
	//init logger
	if err := logger.InitLogger("info"); err != nil {
		fmt.Println(err)
	}

	app.Run(logger.Log)

}
