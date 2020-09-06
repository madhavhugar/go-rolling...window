package main

import "time"

// RollingWindowDataStructure represents interface for
// any data structure that wants to implement a rolling window
type RollingWindowDataStructure interface {
	AppendRequest(time.Time) // should add incoming request to the window
	GetCount() int           // should return the current window size
}
