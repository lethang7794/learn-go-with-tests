package poker

import (
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	t.Run("win call once & record correct winner", func(t *testing.T) {
		in := strings.NewReader("Andy wins")
		store := &StubPlayerStore{}
		cli := &CLI{store, in}

		cli.PlayPoker()

		assertWinner(t, store, "Andy")
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
