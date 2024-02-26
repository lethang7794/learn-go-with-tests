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

		if out.String() != userPrompt {
			t.Errorf("got %#v, want %#v", out.String(), userPrompt)
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

		wantPrompt := userPrompt + "(You're so silly)"
		if out.String() != wantPrompt {
			t.Errorf("got %#v, want %#v", out.String(), wantPrompt)
		}
	})
}
