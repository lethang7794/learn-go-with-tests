package context_aware_reader

import (
	"log"
	"strings"
	"testing"
)

func TestContextAwareReader(t *testing.T) {
	t.Run("a normal reader", func(t *testing.T) {
		reader := strings.NewReader("123456")
		buf := make([]byte, 3)

		n, err := reader.Read(buf)
		if err != nil {
			log.Fatalf("could not read from buf: %v", err)
		}
		assertBufRead(t, string(buf), n, "123")
		n, err = reader.Read(buf)
		if err != nil {
			log.Fatalf("could not read from buf: %v", err)
		}
		assertBufRead(t, string(buf), n, "456")
	})

	t.Run("acts as a normal reader", func(t *testing.T) {
		reader := strings.NewReader("123456")
		reader = NewCancellableReader(reader)
		buf := make([]byte, 3)

		n, err := reader.Read(buf)
		if err != nil {
			log.Fatalf("could not read from buf: %v", err)
		}
		assertBufRead(t, string(buf), n, "123")
		n, err = reader.Read(buf)
		if err != nil {
			log.Fatalf("could not read from buf: %v", err)
		}
		assertBufRead(t, string(buf), n, "456")
	})

	t.Run("doesn't read after cancelled", func(t *testing.T) {
		reader := strings.NewReader("123456")
		ctx, cancel := context.WithCancel(context.Background())
		reader = NewCancellableReader(ctx, reader)
		buf := make([]byte, 3)

		n, err := reader.Read(buf)
		if err != nil {
			log.Fatalf("could not read from buf: %v", err)
		}
		assertBufRead(t, string(buf), n, "123")

		cancel()

		n, err = reader.Read(buf)
		if err == nil {
			t.Fatalf("expected an error, but didn't get one")
		}
		if n != 0 {
			t.Errorf("expected 0 bytes to be read, got %v bytes read", n)
		}
	})
}

func assertBufRead(t *testing.T, got string, n int, want string) {
	t.Helper()
	if got != want {
		t.Errorf("wrong bytes read: got %q, want %q", got, want)
	}
	if n != len(want) {
		t.Errorf("wrong number of bytes read: got %v, want %v", n, len(want))
	}
}
