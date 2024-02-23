package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type League []Player

func NewLeague(reader io.Reader) (League, error) {
	var league League
	err := json.NewDecoder(reader).Decode(&league)
	if err != nil {
		return nil, fmt.Errorf("failed parsing league from JSON: %v", err)
	}
	return league, nil
}

func (l League) Find(name string) *Player {
	for i, player := range l {
		if player.Name == name {
			return &l[i]
		}
	}
	return nil
}
