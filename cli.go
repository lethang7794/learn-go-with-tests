package poker

import (
	"bufio"
	"io"
	"strings"
)

type CLI struct {
	store   PlayerStore
	scanner *bufio.Scanner
}

func NewCLI(store PlayerStore, in io.Reader) *CLI {
	return &CLI{
		store:   store,
		scanner: bufio.NewScanner(in),
	}
}

func (c *CLI) PlayPoker() {
	c.scanner.Scan()
	line := c.scanner.Text()
	winner := extractWinner(line)
	c.store.RecordWin(winner)
}

func extractWinner(line string) string {
	winner := strings.Replace(line, " wins", "", 1)
	return winner
}
