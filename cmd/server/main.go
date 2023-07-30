// Сервер для приема, обработки и хранения метрик
// Разработано по техническому заданию на курсе Golang Developer
// Автор: Денис Дружинин, h2p2f
//
// Server for receiving, processing and storing metrics
// Developed according to the technical task in the Golang Developer course
// Author: Denis Druzhinin, h2p2f

package main

import (
	"fmt"
	"github.com/h2p2f/practicum-metrics/internal/server/app"
) //nolint:typecheck

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

// Запуск сервера
// Server start
func main() {
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)
	app.Run()

}
