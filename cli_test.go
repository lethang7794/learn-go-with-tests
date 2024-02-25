package poker

import (
	"fmt"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	t.Run("win call once & record correct winner", func(t *testing.T) {
		winner := "Andy"
		input := fmt.Sprintf("%s wins", winner)
		in := strings.NewReader(input)
		store := &StubPlayerStore{}
		cli := &CLI{store, in}

		cli.PlayPoker()

		if len(store.winCalls) != 1 {
			t.Errorf("want win to be called once, but it's called %v times (%v)", len(store.winCalls), store.winCalls)
		}

		if store.winCalls[0] != winner {
			t.Errorf("got %#v, want %#v", store.winCalls[0], winner)
		}
	})
}
