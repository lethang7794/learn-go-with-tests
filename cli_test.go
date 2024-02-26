package poker

import (
	"bytes"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	t.Run("prompt the user to enter the number of players & starts the game", func(t *testing.T) {
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

	t.Run("when users enter non numeric value, it print out an error & doesn't start the game", func(t *testing.T) {
		in := strings.NewReader("Not a number\n")
		out := &bytes.Buffer{}
		gameSpy := &GameSpy{}
		cli := NewCLI(in, out, gameSpy)

		cli.PlayPoker()

		if gameSpy.Started == true {
			t.Errorf("Expected game not started, but it does")
		}
	})
}
