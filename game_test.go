package poker

import (
	"fmt"
	"testing"
	"time"
)

func TestGame_StartGame(t *testing.T) {
	t.Run("schedules alert on game start for 5 players", func(t *testing.T) {
		store := &StubPlayerStore{}
		alerter := &SpyBlindAlerter{}
		game := NewGame(store, alerter)

		game.StartGame(5)

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
		game := NewGame(store, alerter)

		game.StartGame(7)

		tt := []scheduledAlert{
			{0 * time.Minute, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		assertScheduledAlerts(t, tt, alerter)
	})
}

func TestGame_FinishGame(t *testing.T) {
	t.Run("Record Andy wins", func(t *testing.T) {
		store := &StubPlayerStore{}
		alerter := &SpyBlindAlerter{}
		game := NewGame(store, alerter)

		game.FinishGame("Andy")

		assertWinner(t, store, "Andy")
	})
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
