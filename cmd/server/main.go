// Server for receiving, processing and storing metrics
// Developed according to the technical task in the Golang Developer course
// Author: Denis Druzhinin, h2p2f

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/h2p2f/practicum-metrics/internal/server/app"
) //nolint:typecheck

// variables for storing the version, date and commit of the build
// set during the build by the command
// go build -ldflags "-X main.buildVersion=1.0.0 -X main.buildDate=2021-09-01 -X main.buildCommit=1234567890"
var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

// Server start
func main() {

	//output information about the version, date and commit of the build
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)

	sigint := make(chan os.Signal, 1)
	connectionsClosed := make(chan struct{})
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	//server start
	app.Run(sigint, connectionsClosed)

	//shutdown signal processing

	<-connectionsClosed

}
