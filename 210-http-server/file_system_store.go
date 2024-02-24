package main

import (
	"encoding/json"
	"os"
	"testing"
)

type FileSystemPlayerStore struct {
	database *json.Encoder
	league   League
}

func NewFileSystemPlayerStore(database *os.File) (*FileSystemPlayerStore, error) {
	database.Seek(0, 0)
	league, err := NewLeague(database)
	if err != nil {
		return nil, err
	}
	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{database}),
		league:   league,
	}, nil
}

func (f *FileSystemPlayerStore) GetLeague() League {
	return f.league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) (int, bool) {
	player := f.league.Find(name)
	if player != nil {
		return player.Score, true
	}
	return 0, false
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.Find(name)
	if player != nil {
		player.Score++
	} else {
		f.league = append(f.league,
			Player{Name: name, Score: 1},
		)
	}
	f.database.Encode(f.league)
}

func createTempFile(t *testing.T, initialData string) (_ *os.File, cleanup func()) {
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
