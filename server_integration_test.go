package poker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIntegrationPlayerHandler(t *testing.T) {
	database, cleanup := createTempFile(t, `[]`)
	defer cleanup()
	store, err := NewFileSystemPlayerStore(database)
	assertNoError(t, err)
	server := NewPlayerServer(store)

	player := "Alpha"
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))

		assertResponseCode(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetLeagueRequest())

		assertResponseCode(t, response.Code, http.StatusOK)

		gotLeague := getLeagueFromResponse(t, response.Body)
		assertLeague(t, gotLeague, League{{
			Name:  player,
			Score: 3,
		}})
	})
}
