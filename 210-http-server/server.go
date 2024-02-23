package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	server := new(PlayerServer)
	server.store = store

	router := http.NewServeMux()

	router.HandleFunc("/players/", server.PlayersHandler)
	router.HandleFunc("/league", server.LeagueHandler)

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

type PlayerStore interface {
	GetPlayerScore(player string) (int, bool)
	RecordWin(name string)
	GetLeague() []Player
}

type Player struct {
	Name  string
	Score int
}
