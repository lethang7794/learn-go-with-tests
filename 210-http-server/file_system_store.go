package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"testing"
)

type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
}

func NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStore {
	return &FileSystemPlayerStore{database: database}
}

func (f FileSystemPlayerStore) GetLeague() League {
	f.database.Seek(0, 0)
	league, err := NewLeague(f.database)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return League{}
		} else {
			log.Fatal(err)
		}
	}
	return league
}

func (f FileSystemPlayerStore) GetPlayerScore(name string) (int, bool) {
	league := f.GetLeague()
	player := league.Find(name)
	if player != nil {
		return player.Score, true
	}
	return 0, false
}

func (f FileSystemPlayerStore) RecordWin(name string) {
	league := f.GetLeague()
	player := league.Find(name)
	if player != nil {
		player.Score++
	} else {
		league = append(league,
			Player{Name: name, Score: 1},
		)
	}
	f.database.Seek(0, 0)
	json.NewEncoder(f.database).Encode(league)
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
