package poker

type CLI struct {
	store PlayerStore
}

func (c CLI) PlayPoker() {
	c.store.RecordWin("Bob")
}
