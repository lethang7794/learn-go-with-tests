package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETPlayers(t *testing.T) {
	t.Run("return Pepper score", func(t *testing.T) {
		response := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "players/Pepper", nil)
		GetPlayerHandler(response, request)

		got := response.Body.String()
		want := "20"

		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

	t.Run("return Floyd score", func(t *testing.T) {
		response := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "players/Floyd", nil)
		GetPlayerHandler(response, request)

		got := response.Body.String()
		want := "10"

		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
}
