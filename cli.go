package poker

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const UserPrompt = "Please enter the number of users: "
const BadInputErrMsg = "Bad value received for number of players, please try again with a number"
const BadWinnerErrorMsg = `Bad value received for winner, type "{Name} wins" to record a win`

var ErrInvalidInputWinner = errors.New("invalid input winner")
var ErrNotEnoughInputWinner = errors.New("not enough input for winner")

type CLI struct {
	in   *bufio.Scanner
	out  io.Writer
	game Game
}

func NewCLI(in io.Reader, out io.Writer, game Game) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		game: game,
	}
}

func (c *CLI) PlayPoker() {
	fmt.Fprintf(c.out, UserPrompt)
	line := c.readLine()
	numberOfPlayers, err := strconv.Atoi(line)
	if err != nil {
		fmt.Fprintf(c.out, BadInputErrMsg)
		return
	}

	c.game.Start(numberOfPlayers)
	line = c.readLine()
	winner, err := extractWinner(line)
	if errors.Is(err, ErrInvalidInputWinner) {
		fmt.Fprintf(c.out, BadWinnerErrorMsg)
		return
	}
	c.game.Finish(winner)

}

func (c *CLI) readLine() string {
	c.in.Scan()
	return c.in.Text()
}

func extractWinner(line string) (string, error) {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return "", ErrNotEnoughInputWinner
	}
	re := regexp.MustCompile(`\w+ wins`)
	found := re.Find([]byte(line))
	if found == nil {
		return "", ErrInvalidInputWinner
	}
	winner := strings.Replace(line, " wins", "", 1)
	return winner, nil
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
