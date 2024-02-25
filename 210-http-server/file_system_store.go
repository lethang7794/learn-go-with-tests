package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"testing"
)

type FileSystemPlayerStore struct {
	database *json.Encoder
	league   League
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
	file.Seek(0, 0)

	fileStat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("could not read file info from file %s: %w", file.Name(), err)
	}
	if fileStat.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, 0)
	}
	league, err := NewLeague(file)
	if err != nil {
		return nil, fmt.Errorf("could not read league from file: %w", err)
	}

	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}, nil
}

func (f *FileSystemPlayerStore) GetLeague() League {
	league := f.league
	slices.SortFunc(league, func(a, b Player) int {
		return cmp.Compare(b.Score, a.Score)
	})
	return league
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
