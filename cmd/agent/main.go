// Агент для обработки и отправки метрик на сервер хранения
// Разработано по техническому заданию на курсе Golang Developer
// Автор: Денис Дружинин, h2p2f
//
// Agent for processing and sending metrics to the storage server
// Developed according to the technical task in the Golang Developer course
// Author: Denis Druzhinin, h2p2f
package main

import (
	"fmt"
	"github.com/h2p2f/practicum-metrics/internal/agent/app" //nolint:typecheck
	"github.com/h2p2f/practicum-metrics/internal/logger"
)

// переменные для хранения версии, даты и коммита сборки
// устанавливаются во время сборки командой
// variables for storing the version, date and commit of the build
// set during the build by the command
// go build -ldflags "-X main.buildVersion=1.0.0 -X main.buildDate=2021-09-01 -X main.buildCommit=1234567890"
var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

// Запуск агента
//
// Agent start

func main() {
	if err := logger.InitLogger("info"); err != nil {
		fmt.Println(err)
		return
	}
	//вывод информации о версии, дате и коммите сборки
	//output information about the version, date and commit of the build
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)

	app.Run(logger.Log)
}
