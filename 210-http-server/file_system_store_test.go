package main

import (
	"io"
	"strings"
	"testing"
)

type FileSystemPlayerStore struct {
	database io.ReadSeeker
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
	f.database.Seek(0, 0)
	league, _ := NewLeague(f.database)
	return league
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

		assertLeague(t, got, want)

		// 2nd time
		got = store.GetLeague()
		assertLeague(t, got, want)
	})
}
