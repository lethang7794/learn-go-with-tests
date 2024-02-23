package main

import (
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	database, cleanup := createTempFile(t, `[
	{ "Name": "Alpha", "Score": 10 },
	{ "Name": "Beta", "Score": 20 }
]`)
	defer cleanup()
	store := NewFileSystemPlayerStore(database)

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

	t.Run("store wins for existing player", func(t *testing.T) {
		store.RecordWin("Beta")

		got, _ := store.GetPlayerScore("Beta")
		want := 21
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

	t.Run("store wins for new player", func(t *testing.T) {
		store.RecordWin("Gamma")

		got, _ := store.GetPlayerScore("Gamma")
		want := 1
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
}
