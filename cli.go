package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const userPrompt = "Please enter the number of users: "

type CLI struct {
	in   *bufio.Scanner
	out  io.Writer
	game Game
}

type Game struct {
	store   PlayerStore
	alerter BlindAlerter
}

func NewCLI(store PlayerStore, in io.Reader, out io.Writer, alerter BlindAlerter) *CLI {
	return &CLI{
		in:  bufio.NewScanner(in),
		out: out,
		game: Game{
			store:   store,
			alerter: alerter,
		},
	}
}

func (c *CLI) PlayPoker() {
	fmt.Fprintf(c.out, userPrompt)
	line := c.readLine()
	numberOfPlayers, _ := strconv.Atoi(line)
	c.scheduleBlindAlerts(numberOfPlayers)
	line = c.readLine()
	winner := extractWinner(line)
	c.game.store.RecordWin(winner)
}

func (c *CLI) scheduleBlindAlerts(numberOfPlayers int) {
	const baseTime = 5
	blindIncrement := time.Duration(baseTime+numberOfPlayers) * time.Minute
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		c.game.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime += blindIncrement
	}
}

func (c *CLI) readLine() string {
	c.in.Scan()
	return c.in.Text()
}

func extractWinner(line string) string {
	winner := strings.Replace(line, " wins", "", 1)
	return winner
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
