package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTeapotHandler_Acceptance(t *testing.T) {
	t.Run("returns teapot status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(TeapotHandler))

		response, err := http.Get(server.URL)
		if err != nil {
			t.Fatal(err)
		}

		got := response.StatusCode
		want := http.StatusTeapot
		if got != want {
			t.Errorf("wrong status: got %#v, want %#v", got, want)
		}
	})
}
