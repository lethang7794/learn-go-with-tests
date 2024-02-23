package main

import (
	"io"
	"os"
	"testing"
)

type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
}

func (f FileSystemPlayerStore) GetPlayerScore(player string) (int, bool) {
	for _, p := range f.GetLeague() {
		if p.Name == player {
			return p.Score, true
		}
	}
	return 0, false
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
	database, cleanup := createTempFile(t, `[
	{ "Name": "Alpha", "Score": 10 },
	{ "Name": "Beta", "Score": 20 }
]`)
	defer cleanup()
	store := FileSystemPlayerStore{database}

	t.Run("league from a reader", func(t *testing.T) {
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
	
	t.Run("get player score", func(t *testing.T) {
		got, _ := store.GetPlayerScore("Beta")
		want := 20
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
}

func createTempFile(t *testing.T, initialData string) (_ io.ReadWriteSeeker, cleanup func()) {
	temp, err := os.CreateTemp("", "db")
	if err != nil {
		t.Fatalf("failed creating temp file: %v", err)
	}
	_, err = temp.Write([]byte(initialData))
	if err != nil {
		t.Fatalf("failed writing initial data: %v", err)
	}
	return temp, func() {
		temp.Close()
		os.Remove(temp.Name())
	}
}
