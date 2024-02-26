package poker

import (
	"bufio"
	"io"
	"strings"
	"time"
)

type CLI struct {
	store   PlayerStore
	scanner *bufio.Scanner
}

func NewCLI(store PlayerStore, in io.Reader, alerter Alerter) *CLI {
	return &CLI{
		store:   store,
		scanner: bufio.NewScanner(in),
	}
}

func (c *CLI) PlayPoker() {
	line := c.readLine()
	winner := extractWinner(line)
	c.store.RecordWin(winner)
}

func (c *CLI) readLine() string {
	c.scanner.Scan()
	return c.scanner.Text()
}

func extractWinner(line string) string {
	winner := strings.Replace(line, " wins", "", 1)
	return winner
}

type Alerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}

type SpyBlindAlerter struct {
	alerts []struct {
		scheduledAt time.Duration
		amount      int
	}
}

func (a SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	a.alerts = append(a.alerts, struct {
		scheduledAt time.Duration
		amount      int
	}{scheduledAt: duration, amount: amount})
}
