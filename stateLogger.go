package main

import "time"

// StateStorage represents the interface to
// read and write state of the rolling window
type StateStorage interface {
	Writer(time.Time)
	Reader() RollingWindowDataStructure
}
