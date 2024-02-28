package error_types

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDumbGetter(t *testing.T) {
	t.Run("when you don't get a 200 you get a status error", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusTeapot)
		}))
		defer svr.Close()

		_, err := DumbGetter(svr.URL)

		if err == nil {
			t.Fatal("expected an error")
		}

		got := err
		want := BadStatusError{
			URL:    svr.URL,
			Status: http.StatusTeapot,
		}
		var target BadStatusError
		if !errors.As(got, &target) {
			t.Fatalf("got %#v, want %#v", got, want)
		}
		if target.Status != http.StatusTeapot {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})
}
