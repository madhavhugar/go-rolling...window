package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_GET_CardinalityHandler(t *testing.T) {
	rw := RollingWindow{
		&Window{},
		60,
		LogFile{"test_state.log", 0644},
		time.RFC3339Nano,
	}
	t.Run("returns the latest count", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		handler := cardinality(&rw)
		handler(res, req)

		got := res.Body.String()
		want := "{\"count\":1}\n"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	cleanUp(&rw)
}
