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
	fmt.Fprint(writer, p.store.GetPlayerScore(player))
}

type PlayerStore interface {
	GetPlayerScore(player string) int
}

func GetPlayerScore(player string) string {
	if player == "Floyd" {
		return "10"
	}
	if player == "Pepper" {
		return "20"
	}
	return ""
}

type StubPlayerStore struct {
	scores map[string]int
}

func (s *StubPlayerStore) GetPlayerScore(player string) int {
	return s.scores[player]
}
