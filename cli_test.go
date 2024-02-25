package poker

import "testing"

func TestCLI(t *testing.T) {
	t.Run("win call once", func(t *testing.T) {
		store := &StubPlayerStore{}
		cli := &CLI{store}

		cli.PlayPoker()

		if len(store.winCalls) != 1 {
			t.Errorf("want win to be called once, but it's called %v times (%v)", len(store.winCalls), store.winCalls)
		}
	})
}
