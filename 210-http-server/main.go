package main

import (
	"log"
	"net/http"
)

func main() {
	store := &InMemoryPlayerStore{}
	handler := &PlayerHandler{store}
	err := http.ListenAndServe(":5000", handler)
	if err != nil {
		log.Fatal(err)
	}
}

type InMemoryPlayerStore struct {
}

func (f *InMemoryPlayerStore) GetPlayerScore(player string) (score int, ok bool) {
	return 12345, true
}
