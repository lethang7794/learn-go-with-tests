package main

import (
	"io"
	"testing"
)

func TestTape_Write(t *testing.T) {
	t.Run("override data", func(t *testing.T) {
		initial := "abcde"
		override := "123"

		file, cleanup := createTempFile(t, initial)
		defer cleanup()
		tape := Tape{file}
		tape.Write([]byte(override))

		tape.file.Seek(0, 0)
		got, err := io.ReadAll(tape.file)

		want := override
		if err != nil {
			t.Fatalf("could not read back file: %v", err)
		}
		if string(got) != want {
			t.Errorf("got %#v, want %#v", string(got), want)
		}
	})
}
