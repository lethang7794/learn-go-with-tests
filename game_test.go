package poker

import (
	"fmt"
	"io"
	"testing"
	"time"
)

func TestGame_Start(t *testing.T) {
	t.Run("schedules alert on game start for 5 players", func(t *testing.T) {
		store := &StubPlayerStore{}
		alerter := &SpyBlindAlerter{}
		game := NewTexasHoldem(store, alerter)

		game.Start(5, io.Discard)

		tt := []scheduledAlert{
			{0 * time.Minute, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		assertScheduledAlerts(t, tt, alerter)
	})

	t.Run("schedules alert on game start for 7 players", func(t *testing.T) {
		store := &StubPlayerStore{}
		alerter := &SpyBlindAlerter{}
		game := NewTexasHoldem(store, alerter)

		game.Start(7, io.Discard)

		tt := []scheduledAlert{
			{0 * time.Minute, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		assertScheduledAlerts(t, tt, alerter)
	})
}

func TestGame_Finish(t *testing.T) {
	t.Run("Record Andy wins", func(t *testing.T) {
		store := &StubPlayerStore{}
		alerter := &SpyBlindAlerter{}
		game := NewTexasHoldem(store, alerter)

		game.Finish("Andy")

		assertWinner(t, store, "Andy")
	})
}

func assertScheduledAlert(t *testing.T, alert scheduledAlert, c scheduledAlert) {
	t.Helper()
	if alert.amount != c.amount {
		t.Errorf("got %#v, want %#v", alert.amount, c.amount)
	}
	if alert.scheduledAt != c.scheduledAt {
		t.Errorf("got %#v, want %#v", alert.scheduledAt, c.scheduledAt)
	}
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

func assertScheduledAlerts(t *testing.T, tt []scheduledAlert, alerter *SpyBlindAlerter) {
	for i, c := range tt {
		name := fmt.Sprintf("%v schedule for %v", c.amount, c.scheduledAt)
		t.Run(name, func(t *testing.T) {
			if len(alerter.alerts) <= i {
				t.Fatalf("alert %v (%v) was not scheduled: %v", i, c.amount, alerter.alerts)
			}
			alert := alerter.alerts[i]
			assertScheduledAlert(t, alert, c)
		})
	}
}
