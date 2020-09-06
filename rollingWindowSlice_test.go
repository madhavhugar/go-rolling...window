package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_RollingWindowSlice_Get(t *testing.T) {
	now := time.Now()
	rws := RollingWindowSlice{
		window:       &[]time.Time{now},
		size:         60,
		currentIndex: 0,
	}

	assert.Equal(t, rws.get(), now)
}

func Test_RollingWindowSlice_HasNext(t *testing.T) {
	now := time.Now()
	// should return true when next element exists in iterator
	rws := RollingWindowSlice{
		window:       &[]time.Time{now, now},
		size:         60,
		currentIndex: 0,
	}
	assert.Equal(t, rws.hasNext(), true)

	// should return false when next element does not exist in iterator
	rws = RollingWindowSlice{
		window:       &[]time.Time{now},
		size:         60,
		currentIndex: 0,
	}
	assert.Equal(t, rws.hasNext(), false)
}

func Test_RollingWindowSlice_Next(t *testing.T) {
	now := time.Now()
	oneMinLater := now.Add(1 * time.Minute)
	twoMinLater := now.Add(2 * time.Minute)
	rws := RollingWindowSlice{
		window: &[]time.Time{
			now,
			oneMinLater,
			twoMinLater,
		},
		size:         60,
		currentIndex: 0,
	}

	// should return the next element if it
	// exists and return nil error
	nextVal, err := rws.next()
	assert.Equal(t, nextVal, oneMinLater)
	assert.Equal(t, err, nil)
	nextVal, err = rws.next()
	assert.Equal(t, nextVal, twoMinLater)
	assert.Equal(t, err, nil)

	// should return error when next element does not exist
	nextVal, err = rws.next()
	assert.NotEqual(t, err, nil)
}

func Test_RollingWindowSlice_Remove(t *testing.T) {
	now := time.Now()
	oneMinLater := now.Add(1 * time.Minute)
	twoMinLater := now.Add(2 * time.Minute)
	rws := RollingWindowSlice{
		window: &[]time.Time{
			now,
			oneMinLater,
			twoMinLater,
		},
		size:         60,
		currentIndex: 0,
	}

	rws.remove()
	assert.Equal(t, len(*rws.window), 2)
	assert.Equal(t, rws.get(), oneMinLater)
}
