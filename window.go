package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// LogFile represents a log file to which the logs are appended
type LogFile struct {
	name        string
	permissions os.FileMode
}

// Window represents a simple data strcuture that stores
// timestamps
type Window []time.Time

// RollingWindow represents the rolling time window and
// relevant metadata
type RollingWindow struct {
	*Window            // embedded type representing the window data structure
	size       int     // time window size in seconds
	logFile    LogFile // log file metadata to store rolling window logs
	timeFormat string  // time format to parse the timestamps
}

// loadRollingWindowFromFile populates the RollingWindow
// with the relevant timestamps from the log file
func loadRollingWindowFromFile(rw *RollingWindow) {
	wb := rw.computeWindowBoundary(time.Now())
	file, err := os.OpenFile(
		rw.logFile.name,
		os.O_RDONLY,
		rw.logFile.permissions,
	)
	if err != nil && strings.Contains(err.Error(), "no such file or directory") {
		return
	}
	handleOnError(err, "on opening file at startup")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		timeLog, err := time.Parse(
			rw.timeFormat,
			scanner.Text(),
		)
		handleOnError(err, "on parsing time from file")
		if isTimeStampWithinWindow(wb, timeLog) {
			rw.addRequestToWindow(timeLog)
		}
	}

	if err := scanner.Err(); err != nil {
		handleOnError(err, "on scanning file at startup")
	}
}

// isTimeStampWithinWindow checks if a given
// timestmap is within a given window boundary
func isTimeStampWithinWindow(
	windowBound time.Time,
	timestamp time.Time,
) bool {
	if timestamp.After(windowBound) {
		return true
	}
	return false
}

// computeWindowBoundary computes the left window edge
// with the right edge as the given timestamp
func (rw *RollingWindow) computeWindowBoundary(
	now time.Time,
) time.Time {
	windowSize := -time.Duration(rw.size)
	return now.Add(windowSize * time.Second)
}

func (rw *RollingWindow) addRequestToWindow(
	timestamp time.Time,
) {
	*rw.Window = append(*rw.Window, timestamp)
}

// garbageCollectOlderTimestamps removes the unrelevant
// timestamps from the Window
func (rw *RollingWindow) garbageCollectOlderTimestamps(
	now time.Time,
) {
	wb := rw.computeWindowBoundary(now)
	var boundIndex int
	var garbageCollect = false
	for index, timestamp := range *rw.Window {
		if !isTimeStampWithinWindow(wb, timestamp) {
			boundIndex = index
			garbageCollect = true
		}
	}
	if garbageCollect {
		*rw.Window = (*rw.Window)[boundIndex+1 : len(*rw.Window)]
	}
}

// appendRequestToLog appends given timestamp
// to the Window
func (rw *RollingWindow) appendRequestToLog(
	timestamp time.Time,
) {
	f, err := os.OpenFile(
		rw.logFile.name,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		rw.logFile.permissions,
	)
	handleOnError(err, "on opening file to write")
	defer f.Close()

	logStatement := fmt.Sprintf(
		"%s\n",
		timestamp.Format(rw.timeFormat),
	)
	_, err = f.WriteString(logStatement)
	handleOnError(err, "on writing to log file")
}
