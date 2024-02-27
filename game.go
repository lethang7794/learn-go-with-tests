package poker

import (
	"io"
	"time"
)

type Game interface {
	Start(numberOfPlayers int, alertDestination io.Writer)
	Finish(winner string)
}

type TexasHoldem struct {
	store   PlayerStore
	alerter BlindAlerter
}

func NewTexasHoldem(store PlayerStore, alerter BlindAlerter) *TexasHoldem {
	return &TexasHoldem{store: store, alerter: alerter}
}

func (g TexasHoldem) Start(numberOfPlayers int, alertDestiation io.Writer) {
	const baseTime = 5
	blindIncrement := time.Duration(baseTime+numberOfPlayers) * time.Minute
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		g.alerter.ScheduleAlertAt(blindTime, blind, alertDestiation)
		blindTime += blindIncrement
	}
}

func (g TexasHoldem) Finish(winner string) {
	g.store.RecordWin(winner)
}

type GameSpy struct {
	StartCalled     bool
	StartCalledWith int

	FinishCalledWith string
}

func (g *GameSpy) Start(numberOfPlayers int, to io.Writer) {
	g.StartCalledWith = numberOfPlayers
	g.StartCalled = true
}

func (g *GameSpy) Finish(winner string) {
	g.FinishCalledWith = winner
}
