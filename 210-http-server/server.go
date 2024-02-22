package main

import (
	"fmt"
	"net/http"
	"strings"
)

type PlayerServer struct {
	store  PlayerStore
	router *http.ServeMux
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	server := &PlayerServer{
		store:  store,
		router: http.NewServeMux(),
	}

	server.router.HandleFunc("/players/", server.PlayersHandler)
	server.router.HandleFunc("/league", server.LeagueHandler)

	return server
}

func (p *PlayerServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	p.router.ServeHTTP(writer, request)
}

func (p *PlayerServer) LeagueHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
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
}
