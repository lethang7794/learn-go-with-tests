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

	t.Run("it schedules printing of blinding values", func(t *testing.T) {
		in := strings.NewReader("Charlie wins")
		store := &StubPlayerStore{}
		spyBlindAlerter := &SpyBlindAlerter{}
		cli := NewCLI(store, in, spyBlindAlerter)

		cli.PlayPoker()

		if len(spyBlindAlerter.alerts) != 1 {
			t.Errorf("expected a blind alert to be scheduled")
		}
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
