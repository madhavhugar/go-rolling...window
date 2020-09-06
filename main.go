package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var rollingWin RollingWindow
var mutex sync.Mutex
var sliceRollingWindow *RollingWindowSlice

func main() {
	logFile := LogFile{
		name:        "state.log",
		permissions: 0644,
	}
	rollingWin = RollingWindow{
		Window:     &Window{},
		size:       60,
		timeFormat: time.RFC3339Nano,
		logFile:    logFile,
	}
	loadRollingWindowFromFile(&rollingWin)

	sliceRollingWindow = NewRollingWindowSlice(60)

	mux := http.NewServeMux()
	mux.HandleFunc("/", cardinality(&rollingWin))
	mux.HandleFunc("/cardinality", cardinalityUsingAbstractedSlice(sliceRollingWindow))

	serverAddress := "0.0.0.0:9090"
	fmt.Println("Starting server at", serverAddress)
	err := http.ListenAndServe(serverAddress, mux)
	handleOnError(err, "on starting the server")
}
