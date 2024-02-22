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
	score, ok := p.store.GetPlayerScore(player)
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(writer, score)
}

type PlayerStore interface {
	GetPlayerScore(player string) (int, bool)
}

type StubPlayerStore struct {
	scores map[string]int
}

func (s *StubPlayerStore) GetPlayerScore(player string) (score int, ok bool) {
	score, ok = s.scores[player]
	if ok {
		return score, true
	}
	return score, false
}
