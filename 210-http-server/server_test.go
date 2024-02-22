package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETPlayers(t *testing.T) {
	store := &StubPlayerStore{
		scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
			"Beta":   0,
		},
	}
	server := NewPlayerServer(store)

	t.Run("return Pepper score", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newGetScoreRequest("Pepper")

		server.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("return Floyd score", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newGetScoreRequest("Floyd")

		server.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("return 404 on missing players", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newGetScoreRequest("Alpha")

		server.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusNotFound)
	})

	t.Run("return 200 on player with 0 score", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newGetScoreRequest("Beta")

		server.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "0")
	})

}

func TestStoreWins(t *testing.T) {
	store := &StubPlayerStore{
		scores: map[string]int{},
	}
	server := NewPlayerServer(store)

	t.Run("record win when POST a player", func(t *testing.T) {
		player := "Alpha"
		response := httptest.NewRecorder()
		request := newPostWinRequest(player)

		server.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusAccepted)

		if got, want := len(store.winCalls), 1; got != want {
			t.Errorf("wrong number of call to RecordWin: got %#v, want %#v", got, want)
		}
		if got, want := store.winCalls[0], player; got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
}

func TestLeague(t *testing.T) {
	t.Run("return 200 on /league", func(t *testing.T) {
		store := &StubPlayerStore{}
		server := NewPlayerServer(store)
		response := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)

		server.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusOK)

		var got []Player
		err := json.NewDecoder(response.Body).Decode(&got)
		if err != nil {
			t.Errorf("failed parsing JSON from %q: error %v", response.Body, err)
		}
	})
}

func newGetScoreRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/players/"+name, nil)
	return request
}

func newPostWinRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, "/players/"+name, nil)
	return request
}

func assertResponseBody(t *testing.T, got string, want string) {
	t.Helper()
	if got != want {
		t.Errorf("wrong response body: got %#v, want %#v", got, want)
	}
}

func assertResponseCode(t *testing.T, got int, want int) {
	t.Helper()
	if got != want {
		t.Errorf("invalid code: got %#v, want %#v", got, want)
	}
}
