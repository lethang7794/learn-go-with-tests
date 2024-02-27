package poker

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

const retryTime = 500 * time.Millisecond

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

	t.Run("start game with 3 players and finish game with 'Chris' as winner", func(t *testing.T) {
		in := strings.NewReader("3\nChris wins")
		out := &bytes.Buffer{}
		gameSpy := &GameSpy{}
		cli := NewCLI(in, out, gameSpy)

		cli.PlayPoker()

		assertGameStartWith(t, gameSpy, 3)
		assertGameFinishWith(t, gameSpy, "Chris")
	})

	t.Run("start game with 8 players and finish game with 'James' as winner", func(t *testing.T) {
		in := strings.NewReader("8\nJames wins")
		out := &bytes.Buffer{}
		gameSpy := &GameSpy{}
		cli := NewCLI(in, out, gameSpy)

		cli.PlayPoker()

		assertGameStartWith(t, gameSpy, 8)
		assertGameFinishWith(t, gameSpy, "James")
	})

	t.Run(`when user not enter "{Name} wins", show an error message`, func(t *testing.T) {
		in := strings.NewReader("8\nYou're so silly")
		out := &bytes.Buffer{}
		gameSpy := &GameSpy{}
		cli := NewCLI(in, out, gameSpy)

		cli.PlayPoker()

		assertGameStartWith(t, gameSpy, 8)
		assertMessageSendToUser(t, out, UserPrompt, BadWinnerErrorMsg)
	})
}

func assertGameFinishWith(t *testing.T, gameSpy *GameSpy, winner string) {
	t.Helper()

	passed := retryUntil(retryTime, func() bool {
		return gameSpy.FinishCalledWith != winner
	})

	if passed {
		t.Errorf("got %#v, winner %#v", gameSpy.FinishCalledWith, winner)
	}
}

func assertGameNotStarted(t *testing.T, gameSpy *GameSpy) {
	t.Helper()
	if gameSpy.StartCalled == true {
		t.Errorf("game should have not started")
	}
}

func assertMessageSendToUser(t *testing.T, out *bytes.Buffer, messages ...string) {
	t.Helper()
	message := strings.Join(messages, "")
	if out.String() != message {
		t.Errorf("got %#v, want %#v", out.String(), message)
	}
}

func assertGameStartWith(t *testing.T, gameSpy *GameSpy, numberOfPlayers int) {
	t.Helper()
	if gameSpy.StartCalledWith != numberOfPlayers {
		t.Errorf("got %#v, want %#v", gameSpy.StartCalledWith, numberOfPlayers)
	}
}

func retryUntil(d time.Duration, fn func() bool) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if fn() {
			return true
		}
	}
	return false
}
