package main

import (
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
	handler := &PlayerHandler{store}

	t.Run("return Pepper score", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newGetScoreRequest("Pepper")

		handler.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("return Floyd score", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newGetScoreRequest("Floyd")

		handler.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("return 404 on missing players", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newGetScoreRequest("Alpha")

		handler.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusNotFound)
	})

	t.Run("return 200 on player with 0 score", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newGetScoreRequest("Beta")

		handler.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "0")
	})

}

func TestStoreWins(t *testing.T) {
	store := &StubPlayerStore{
		scores: map[string]int{},
	}
	handler := &PlayerHandler{store}

	t.Run("record win when POST a player", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := newPostWinRequest("Alpha")

		handler.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusAccepted)

		got := len(store.winCalls)
		want := 1
		if got != want {
			t.Errorf("wrong number of call to RecordWin: got %#v, want %#v", got, want)
		}
	})
}

func newGetScoreRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "players/"+name, nil)
	return request
}

func newPostWinRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, "players/"+name, nil)
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
