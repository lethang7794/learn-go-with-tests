package main

import "sync"

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

func (s *InMemoryPlayerStore) GetLeague() League {
	league := League{}
	for key, val := range s.scores {
		league = append(league, Player{
			Name:  key,
			Score: val,
		})
	}
	return league
}
