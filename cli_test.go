package poker

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCLI(t *testing.T) {
	t.Run("Record Andy wins", func(t *testing.T) {
		in := strings.NewReader("5\nAndy wins")
		store := &StubPlayerStore{}
		alerter := &SpyBlindAlerter{}
		cli := NewCLI(store, in, os.Stdout, alerter)

		cli.PlayPoker()

		assertWinner(t, store, "Andy")
	})

	t.Run("Record Bob wins", func(t *testing.T) {
		in := strings.NewReader("5\nBob wins")
		store := &StubPlayerStore{}
		alerter := &SpyBlindAlerter{}
		cli := NewCLI(store, in, os.Stdout, alerter)

		cli.PlayPoker()

		assertWinner(t, store, "Bob")
	})

	t.Run("it schedules printing of blinding values", func(t *testing.T) {
		in := strings.NewReader("5\nCharlie wins")
		store := &StubPlayerStore{}
		spyBlindAlerter := &SpyBlindAlerter{}
		cli := NewCLI(store, in, os.Stdout, spyBlindAlerter)

		cli.PlayPoker()

		cases := []scheduledAlert{
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
		for i, c := range cases {
			name := fmt.Sprintf("%v schedule for %v", c.amount, c.scheduledAt)
			t.Run(name, func(t *testing.T) {
				if len(spyBlindAlerter.alerts) <= i {
					t.Fatalf("alert %v (%v) was not scheduled: %v", i, c.amount, spyBlindAlerter.alerts)
				}

				alert := spyBlindAlerter.alerts[i]

				assertScheduledAlert(t, alert, c)
			})
		}
	})

	t.Run("its prompt the user to enter the number of players", func(t *testing.T) {
		store := &StubPlayerStore{}
		in := strings.NewReader("7\n")
		out := &bytes.Buffer{}
		alerter := &SpyBlindAlerter{}
		cli := NewCLI(store, in, out, alerter)

		cli.PlayPoker()

		got := out.String()
		want := userPrompt
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}

		tt := []scheduledAlert{
			{0 * time.Minute, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

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
