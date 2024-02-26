package poker

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

type CLI struct {
	store   PlayerStore
	scanner *bufio.Scanner
	alerter Alerter
}

func NewCLI(store PlayerStore, in io.Reader, alerter Alerter) *CLI {
	return &CLI{
		store:   store,
		scanner: bufio.NewScanner(in),
		alerter: alerter,
	}
}

func (c *CLI) PlayPoker() {
	c.scheduleBlindAlerts()

	line := c.readLine()
	winner := extractWinner(line)
	c.store.RecordWin(winner)
}

func (c *CLI) scheduleBlindAlerts() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		c.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime += 10 * time.Minute
	}
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

type scheduledAlert struct {
	scheduledAt time.Duration
	amount      int
}

func (a scheduledAlert) String() string {
	return fmt.Sprintf("%v chips at %v", a.amount, a.scheduledAt)
}

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

func (a *SpyBlindAlerter) ScheduleAlertAt(duration time.Duration, amount int) {
	a.alerts = append(a.alerts, scheduledAlert{
		scheduledAt: duration, amount: amount,
	})
}
