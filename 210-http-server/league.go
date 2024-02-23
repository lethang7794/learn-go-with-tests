package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func NewLeague(reader io.Reader) ([]Player, error) {
	var league []Player
	err := json.NewDecoder(reader).Decode(&league)
	if err != nil {
		return nil, fmt.Errorf("failed parsing league from JSON: %v", err)
	}
	return league, nil
}
