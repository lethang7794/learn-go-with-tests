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

		assertMessageSendToUser(t, out, UserPrompt)
		assertGameStartWith(t, gameSpy, 7)
	})

	t.Run("when users enter non numeric value, it print out an error & doesn't start the game", func(t *testing.T) {
		in := strings.NewReader("Not a number\n")
		out := &bytes.Buffer{}
		gameSpy := &GameSpy{}
		cli := NewCLI(in, out, gameSpy)

		cli.PlayPoker()

		assertGameNotStarted(t, gameSpy)
		assertMessageSendToUser(t, out, UserPrompt+BadInputErrMsg)
	})
}

func assertGameNotStarted(t *testing.T, gameSpy *GameSpy) {
	t.Helper()
	if gameSpy.Started == true {
		t.Errorf("game should have not started")
	}
}

func assertMessageSendToUser(t *testing.T, out *bytes.Buffer, msgs ...string) {
	t.Helper()
	message := strings.Join(msgs, "")
	if out.String() != message {
		t.Errorf("got %#v, want %#v", out.String(), message)
	}
}

func assertGameStartWith(t *testing.T, gameSpy *GameSpy, numberOfPlayers int) {
	t.Helper()
	if gameSpy.StartedWith != numberOfPlayers {
		t.Errorf("got %#v, want %#v", gameSpy.StartedWith, numberOfPlayers)
	}
}
