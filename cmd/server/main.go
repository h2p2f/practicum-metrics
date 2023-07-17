package main

import (
	"github.com/h2p2f/practicum-metrics/internal/server/app" //nolint:typecheck

	_ "net/http/pprof"
)

//nolint:typecheck

func main() {

	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()
	app.Run()

}
