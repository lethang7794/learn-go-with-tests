package context_aware_reader

import "strings"

func NewCancellableReader(reader *strings.Reader) *strings.Reader {
	return reader
}
