package main

import (
	"bufio"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setUp(now time.Time) RollingWindow {
	window := Window{
		now.Add(-10 * time.Minute),
		now.Add(-5 * time.Minute),
		now.Add(-65 * time.Second),
		now.Add(-55 * time.Second),
		now.Add(-10 * time.Second),
	}
	rollingWin := RollingWindow{
		&window,
		60,
		LogFile{"test_state.log", 0644},
		time.RFC3339Nano,
	}
	return rollingWin
}

func cleanUp(rw *RollingWindow) {
	os.Remove(rw.logFile.name)
}

func Test_garbageCollectOlderTimestamps(t *testing.T) {
	now := time.Now()
	rollingWin := setUp(now)

	expected := Window{
		now.Add(-55 * time.Second),
		now.Add(-10 * time.Second),
	}

	rollingWin.garbageCollectOlderTimestamps(now)
	got := *rollingWin.Window
	assert.Equal(t, expected, got)

	cleanUp(&rollingWin)
}

// func Test_getLatestRequest(t *testing.T) {
// 	now := time.Now()
// 	rollingWin := setUp(now)
// 	expected := now.Add(-10 * time.Second)

// 	got := rollingWin.getLatestRequest()
// 	assert.Equal(t, expected, got)
// }

func Test_isTimeStampWithinWindow(t *testing.T) {
	now := time.Now()
	windowBound := now.Add(-60 * time.Second)
	timeOutsideWindow := now.Add(-65 * time.Second)
	timeWithinWindow := now.Add(-10 * time.Second)

	// should return false when the time is outside the window frame
	got := isTimeStampWithinWindow(windowBound, timeOutsideWindow)
	expected := false
	assert.Equal(t, expected, got)

	// should return true when the time is within the window frame
	got = isTimeStampWithinWindow(windowBound, timeWithinWindow)
	expected = true
	assert.Equal(t, expected, got)
}

func Test_loadRollingWindowFromFile(t *testing.T) {
	// should return an empty RollingWindow object when the log file is missing
	now := time.Now()
	rw := RollingWindow{
		&Window{},
		60,
		LogFile{"test_state.log", 0644},
		time.RFC3339Nano,
	}

	loadRollingWindowFromFile(&rw)
	assert.Equal(t, 0, len(*rw.Window))
	cleanUp(&rw)

	// should populate rolling window with details from the specified log file
	rw = RollingWindow{
		&Window{},
		60,
		LogFile{"test_state.log", 0644},
		time.RFC3339Nano,
	}
	rw.appendRequestToLog(now)
	rollingWin := RollingWindow{
		&Window{},
		60,
		LogFile{"test_state.log", 0644},
		time.RFC3339Nano,
	}
	loadRollingWindowFromFile(&rollingWin)
	assert.Equal(
		t,
		now.Format(rw.timeFormat),
		(*rollingWin.Window)[0].Format(rw.timeFormat),
	)
	cleanUp(&rw)
}

func Test_appendRequestToLog(t *testing.T) {
	// should append timestamp to the specified log file
	now := time.Now()
	rollingWin := RollingWindow{
		&Window{},
		60,
		LogFile{"test_state.log", 0644},
		time.RFC3339Nano,
	}
	rollingWin.appendRequestToLog(now)

	file, _ := os.Open(rollingWin.logFile.name)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		log := scanner.Text()
		expected := now
		got, _ := time.Parse(rollingWin.timeFormat, log)
		assert.Equal(
			t,
			expected.Format(rollingWin.timeFormat),
			got.Format(rollingWin.timeFormat),
		)
	}
	cleanUp(&rollingWin)
}

func Test_computeWindowBoundary(t *testing.T) {
	now := time.Now()
	rollingWin := RollingWindow{
		&Window{},
		60,
		LogFile{"test_state.log", 0644},
		time.RFC3339Nano,
	}

	got := rollingWin.computeWindowBoundary(now)
	expected := now.Add(-60 * time.Second)
	assert.Equal(t, expected, got)
}

func Test_addRequestToWindow(t *testing.T) {
	now := time.Now()
	rollingWin := setUp(now)
	expected := now.Add(-10 * time.Second)

	assert.Equal(
		t,
		expected,
		(*rollingWin.Window)[len(*rollingWin.Window)-1],
	)
}
