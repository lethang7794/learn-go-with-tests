package poker

import (
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		store := &StubPlayerStore{}
		server := NewPlayerServer(store)

		request, _ := http.NewRequest(http.MethodGet, "/game", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusOK)
	})

	t.Run("when we get a message over WebSocket, it's a winner of a game", func(t *testing.T) {
		winner := "Ruth"
		store := &StubPlayerStore{}
		server := NewPlayerServer(store)
		testServer := httptest.NewServer(server.Handler)
		defer testServer.Close()

		wsURl := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/ws"
		conn, _, err := websocket.DefaultDialer.Dial(wsURl, nil)
		if err != nil {
			t.Fatalf("could not open WebSocket connection: %v", err)
		}
		defer conn.Close()

		err = conn.WriteMessage(websocket.TextMessage, []byte(winner))
		if err != nil {
			t.Fatalf("could not send message over WebSocket connection: %v", err)
		}

		assertWinner(t, store, winner)
	})
}

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
	t.Run("return league on /league", func(t *testing.T) {
		wantLeague := League{
			{Name: "First", Score: 1},
			{Name: "Second", Score: 1},
			{Name: "Third", Score: 1},
		}
		store := &StubPlayerStore{league: wantLeague}
		server := NewPlayerServer(store)
		response := httptest.NewRecorder()
		request := newGetLeagueRequest()

		server.ServeHTTP(response, request)

		gotLeague := getLeagueFromResponse(t, response.Body)
		assertLeague(t, gotLeague, wantLeague)
		assertResponseCode(t, response.Code, http.StatusOK)
		assertResponseHeaderContentType(t, response)
	})
}

func getLeagueFromResponse(t *testing.T, body io.Reader) League {
	t.Helper()
	got, _ := NewLeague(body)
	return got
}

func newGetLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
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
		t.Errorf("wrong code: got %#v, want %#v", got, want)
	}
}

func assertLeague(t *testing.T, got League, want League) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func assertResponseHeaderContentType(t *testing.T, response http.ResponseWriter) {
	t.Helper()
	got := response.Header().Get("Content-Type")
	want := jsonContentType
	if got != want {
		t.Errorf("wrong Content-Type: got %#v, want %#v", got, want)
	}
}
