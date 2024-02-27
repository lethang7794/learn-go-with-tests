package poker

import (
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

const timeOut = 10 * time.Millisecond

func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		store := &StubPlayerStore{}
		server := mustMakePlayerServer(t, store, nil)

		request, _ := http.NewRequest(http.MethodGet, "/game", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertResponseCode(t, response.Code, http.StatusOK)
	})

	t.Run("start a game with 3 players, send blinks over websocket connection and declare Ruth the winner", func(t *testing.T) {
		winner := "Ruth"
		wantBlinkAlert := "Blind is 100"

		gameSpy := &GameSpy{BlindAlert: []byte(wantBlinkAlert)}
		store := &StubPlayerStore{}
		server := mustMakePlayerServer(t, store, gameSpy)
		testServer := httptest.NewServer(server)
		defer testServer.Close()

		wsURl := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/ws"
		ws := mustMakeWebSocketConn(t, wsURl)
		defer ws.Close()

		writeWsMessage(t, ws, "3")
		writeWsMessage(t, ws, winner)

		time.Sleep(timeOut) // TODO: remove
		assertGameStartWith(t, gameSpy, 3)
		assertGameFinishWith(t, gameSpy, winner)

		within(t, timeOut, func() {
			assertWebSocketGotMessage(t, ws, wantBlinkAlert)
		})
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
	server := mustMakePlayerServer(t, store, nil)

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
	server := mustMakePlayerServer(t, store, nil)

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
		server := mustMakePlayerServer(t, store, nil)
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

func mustMakePlayerServer(t *testing.T, store PlayerStore, game Game) *PlayerServer {
	t.Helper()
	server, err := NewPlayerServer(store, game)
	if err != nil {
		t.Fatalf("could not create player server: %v", err)
	}
	return server
}

func mustMakeWebSocketConn(t *testing.T, wsURl string) *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial(wsURl, nil)
	if err != nil {
		t.Fatalf("could not open WebSocket connection: %v", err)
	}
	return conn
}

func writeWsMessage(t *testing.T, ws *websocket.Conn, msg string) {
	t.Helper()
	err := ws.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		t.Logf("could not write websocket message: %v", err)
	}
}

func within(t *testing.T, d time.Duration, assert func()) {
	t.Helper()
	done := make(chan bool, 1)
	go func() {
		assert()
		done <- true
	}()
	select {
	case <-time.After(d):
		t.Errorf("time out: %v", d)
	case <-done:
	}

}

func assertWebSocketGotMessage(t *testing.T, ws *websocket.Conn, message string) {
	t.Helper()
	_, msg, _ := ws.ReadMessage()
	if string(msg) != message {
		t.Errorf("got %#v, want %#v", string(msg), message)
	}
}
