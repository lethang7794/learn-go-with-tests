package poker

import (
	"bytes"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	t.Run("its prompt the user to enter the number of players & starts the game", func(t *testing.T) {
		in := strings.NewReader("7\n")
		out := &bytes.Buffer{}
		gameSpy := &GameSpy{}
		cli := NewCLI(in, out, gameSpy)

		cli.PlayPoker()

		got := out.String()
		want := userPrompt
		if got != want {
			t.Errorf("got %#v, want %#v", got, want)
		}

		if gameSpy.StartsWith != 7 {
			t.Errorf("got %#v, want %#v", gameSpy.StartsWith, 7)

		}
	})
}
