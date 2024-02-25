package poker

import "io"

type CLI struct {
	store PlayerStore
	in    io.Reader
}

func (c *CLI) PlayPoker() {
	c.store.RecordWin("Bob")
}
