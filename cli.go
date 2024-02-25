package poker

import (
	"bufio"
	"io"
	"strings"
)

type CLI struct {
	store PlayerStore
	in    io.Reader
}

func NewCLI(store PlayerStore, in io.Reader) *CLI {
	return &CLI{store: store, in: in}
}

func (c *CLI) PlayPoker() {
	reader := bufio.NewScanner(c.in)
	reader.Scan()
	line := reader.Text()
	winner := extractWinner(line)
	c.store.RecordWin(winner)
}

func extractWinner(line string) string {
	winner := strings.Replace(line, " wins", "", 1)
	return winner
}
