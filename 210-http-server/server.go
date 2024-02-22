package main

import (
	"fmt"
	"net/http"
	"strings"
)

type PlayerHandler struct {
	store PlayerStore
}

func (p *PlayerHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	player := strings.TrimPrefix(request.URL.Path, "players/")

	if request.Method == http.MethodGet {
		p.ShowScore(writer, player)
	}
	if request.Method == http.MethodPost {
		p.ProcessWin(writer, player)
	}
}

func (p *PlayerHandler) ShowScore(writer http.ResponseWriter, player string) {
	score, ok := p.store.GetPlayerScore(player)
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(writer, score)
}

func (p *PlayerHandler) ProcessWin(writer http.ResponseWriter, player string) {
	writer.WriteHeader(http.StatusAccepted)
	p.store.RecordWin("Beta")
}

type PlayerStore interface {
	GetPlayerScore(player string) (int, bool)
	RecordWin(name string)
}

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(player string) (score int, ok bool) {
	score, ok = s.scores[player]
	if ok {
		return score, true
	}
	return score, false
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}
