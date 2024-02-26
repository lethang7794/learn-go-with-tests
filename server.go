package poker

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// PlayerStore stores score information about players
type PlayerStore interface {
	GetPlayerScore(player string) (int, bool)
	RecordWin(name string)
	GetLeague() League
}

// Player stores a player with score
type Player struct {
	Name  string
	Score int
}

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

// NewPlayerServer creates a PlayerServer with routing configured
func NewPlayerServer(store PlayerStore) *PlayerServer {
	server := new(PlayerServer)
	server.store = store

	router := http.NewServeMux()

	router.HandleFunc("/players/", server.PlayersHandler)
	router.HandleFunc("/league", server.LeagueHandler)
	router.HandleFunc("/game", server.GameHandler)

	server.Handler = router

	return server
}

const jsonContentType = "application/json"

func (p *PlayerServer) LeagueHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", jsonContentType)
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(p.store.GetLeague())
}

func (p *PlayerServer) PlayersHandler(writer http.ResponseWriter, request *http.Request) {
	player := strings.TrimPrefix(request.URL.Path, "/players/")

	if request.Method == http.MethodGet {
		p.ShowScore(writer, player)
	}
	if request.Method == http.MethodPost {
		p.ProcessWin(writer, player)
	}
}

func (p *PlayerServer) ShowScore(writer http.ResponseWriter, player string) {
	score, ok := p.store.GetPlayerScore(player)
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(writer, score)
}

func (p *PlayerServer) ProcessWin(writer http.ResponseWriter, player string) {
	writer.WriteHeader(http.StatusAccepted)
	p.store.RecordWin(player)
}

func (p *PlayerServer) GameHandler(writer http.ResponseWriter, request *http.Request) {
	files, err := template.ParseFiles("game.html")
	if err != nil {
		http.Error(writer, fmt.Sprintf("could not parse template file: %s", err.Error()), http.StatusInternalServerError)
	}
	files.Execute(writer, nil)
}
