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

func (c *CLI) PlayPoker() {
	reader := bufio.NewScanner(c.in)
	reader.Scan()
	line := reader.Text()
	winner := strings.Replace(line, " wins", "", 1)
	c.store.RecordWin(winner)
}
