package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTeapotHandler(t *testing.T) {
	t.Run("returns teapot status", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := &http.Request{Method: http.MethodGet}

		TeapotHandler(response, request)

		got := response.Code
		want := http.StatusTeapot
		if got != want {
			t.Errorf("wrong status: got %#v, want %#v", got, want)
		}
	})
}
