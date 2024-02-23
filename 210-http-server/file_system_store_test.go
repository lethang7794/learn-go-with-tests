package main

import (
	"encoding/json"
	"io"
	"slices"
	"strings"
	"testing"
)

type FileSystemPlayerStore struct {
	database io.Reader
}

func (f FileSystemPlayerStore) GetPlayerScore(player string) (int, bool) {
	//TODO implement me
	panic("implement me")
}

func (f FileSystemPlayerStore) RecordWin(name string) {
	//TODO implement me
	panic("implement me")
}

func (f FileSystemPlayerStore) GetLeague() []Player {
	league, _ := NewLeague(f.database)
	return league
}

func NewLeague(reader io.Reader) ([]Player, error) {
	var league []Player
	err := json.NewDecoder(reader).Decode(&league)
	if err != nil {
		return nil, err
	}
	return league, nil
}

func TestFileSystemStore(t *testing.T) {
	t.Run("league from a reader", func(t *testing.T) {
		database := strings.NewReader(`[
	{ "Name": "Alpha", "Score": 10 },
	{ "Name": "Beta", "Score": 20 }
]`)
		store := FileSystemPlayerStore{database}

		got := store.GetLeague()
		want := []Player{
			{Name: "Alpha", Score: 10},
			{Name: "Beta", Score: 20},
		}

		if !slices.Equal(got, want) {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
}
