package error_types

import (
	"fmt"
	"io"
	"net/http"
)

type BadStatusError struct {
	URL    string
	Status int
}

func (b BadStatusError) Error() string {
	return fmt.Sprintf("didn't get 200 from %v, got %v", b.URL, b.Status)
}

// DumbGetter will get the string body of url if it gets a 200
func DumbGetter(url string) (string, error) {
	res, err := http.Get(url)
	defer res.Body.Close()

	if err != nil {
		return "", fmt.Errorf("problem fetching from %s, %v", url, err)
	}

	if res.StatusCode != http.StatusOK {
		return "", BadStatusError{
			URL:    url,
			Status: res.StatusCode,
		}
	}

	body, _ := io.ReadAll(res.Body) // ignoring err for brevity
	return string(body), nil
}
