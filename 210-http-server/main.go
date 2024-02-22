package main

import (
	"log"
	"net/http"
	"sync"
)

func main() {
	store := NewInMemoryPlayerStore()
	handler := &PlayerHandler{store}
	err := http.ListenAndServe(":5000", handler)
	if err != nil {
		log.Fatal(err)
	}
}

type InMemoryPlayerStore struct {
	scores map[string]int
	mu     sync.Mutex
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{
		scores: map[string]int{},
	}
}

func (s *InMemoryPlayerStore) GetPlayerScore(player string) (score int, ok bool) {
	score, ok = s.scores[player]
	return score, ok
}

func (s *InMemoryPlayerStore) RecordWin(player string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scores[player]++
}
