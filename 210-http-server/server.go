package main

import (
	"fmt"
	"net/http"
	"strings"
)

type PlayerServer struct {
	store PlayerStore
}

func (p *PlayerServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	mux := http.NewServeMux()

	mux.Handle("/players/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		player := strings.TrimPrefix(request.URL.Path, "/players/")

		if request.Method == http.MethodGet {
			p.ShowScore(writer, player)
		}
		if request.Method == http.MethodPost {
			p.ProcessWin(writer, player)
		}
	}))

	mux.HandleFunc("/league", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	mux.ServeHTTP(writer, request)
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
