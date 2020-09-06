package main

import (
	"errors"
	"time"
)

// RollingWindowSlice is a slice implementation of RollowingWindow
type RollingWindowSlice struct {
	window       *[]time.Time
	size         int
	currentIndex int
}

// NewRollingWindowSlice bla
func NewRollingWindowSlice(size int) *RollingWindowSlice {
	w := &[]time.Time{}
	return &RollingWindowSlice{
		window:       w,
		size:         size,
		currentIndex: 0,
	}
}

func (rws *RollingWindowSlice) get() time.Time {
	return (*rws.window)[rws.currentIndex]
}

func (rws *RollingWindowSlice) hasNext() bool {
	// fmt.Println("len", len(*rws.window), "current ", rws.currentIndex)
	nextIndex := rws.currentIndex + 1
	return len(*rws.window) > nextIndex
}

func (rws *RollingWindowSlice) next() (time.Time, error) {
	if rws.hasNext() {
		rws.currentIndex++
		return (*rws.window)[rws.currentIndex], nil
	}
	return time.Time{}, errors.New("No next element")
}

func (rws *RollingWindowSlice) remove() {
	*rws.window = append((*rws.window)[:rws.currentIndex], (*rws.window)[rws.currentIndex+1:]...)
}

func (rws *RollingWindowSlice) windowSize() int {
	return rws.size
}

func (rws *RollingWindowSlice) GetCount() int {
	return len(*rws.window)
}

func (rws *RollingWindowSlice) AppendRequest(timestamp time.Time) {
	*rws.window = append(*rws.window, timestamp)
}

func computeWindowBoundary(now time.Time, size int) time.Time {
	windowSize := -time.Duration(size)
	return now.Add(windowSize * time.Second)
}

func (rws *RollingWindowSlice) garbageCollect(rw RollingWindowDataStructure) {
	for rws.hasNext() {
		rws.next()
		now := time.Now()
		wb := computeWindowBoundary(now, rws.windowSize())
		for currTimestamp := rws.get(); rws.hasNext(); currTimestamp, _ = rws.next() {
			if !isTimeStampWithinWindow(wb, currTimestamp) {
				rws.remove()
			}
		}
	}
}

// func addRequestToWindow(rw RollingWindowDataStructure, timestamp time.Time) {
// 	rw.AppendRequest(timestamp)
// }
