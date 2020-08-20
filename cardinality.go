package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Response holds the total count of requests
type Response struct {
	Count int `json:"count"`
}

// cardinality handles GET request to route "/"
// and responds with the total number of counts
// within a given window
func cardinality(rw *RollingWindow) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		mutex.Lock()
		now := time.Now()
		rw.garbageCollectOlderTimestamps(now)
		rw.appendRequestToLog(now)
		rw.addRequestToWindow(now)
		mutex.Unlock()

		count := len(*rw.Window)
		fmt.Println("Route => / Method => ", req.Method, "Count => ", count)
		resPayload := Response{count}
		err := json.NewEncoder(res).Encode(&resPayload)
		handleOnError(err, "on encoding response")
	}
}
