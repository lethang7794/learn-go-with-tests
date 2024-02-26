package poker

import "time"

type Game struct {
	store   PlayerStore
	alerter BlindAlerter
}

func NewGame(store PlayerStore, alerter BlindAlerter) *Game {
	return &Game{store: store, alerter: alerter}
}

func (g Game) Start(numberOfPlayers int) {
	const baseTime = 5
	blindIncrement := time.Duration(baseTime+numberOfPlayers) * time.Minute
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		g.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime += blindIncrement
	}
}

func (g Game) Finish(winner string) {
	g.store.RecordWin(winner)
}
