package poker

import (
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	t.Run("Record Andy wins", func(t *testing.T) {
		in := strings.NewReader("Andy wins")
		store := &StubPlayerStore{}
		cli := NewCLI(store, in)

		cli.PlayPoker()

		assertWinner(t, store, "Andy")
	})

	t.Run("Record Bob wins", func(t *testing.T) {
		in := strings.NewReader("Bob wins")
		store := &StubPlayerStore{}
		cli := NewCLI(store, in)

		cli.PlayPoker()

		assertWinner(t, store, "Bob")
	})
}

func assertWinner(t *testing.T, store *StubPlayerStore, winner string) {
	t.Helper()
	if len(store.winCalls) != 1 {
		t.Errorf("winner win to be called once, but it's called %v times (%v)", len(store.winCalls), store.winCalls)
	}
	if store.winCalls[0] != winner {
		t.Errorf("got %#v, winner %#v", store.winCalls[0], winner)
	}
}
