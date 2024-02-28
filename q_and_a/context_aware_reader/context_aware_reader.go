package context_aware_reader

import (
	"context"
	"strings"
)

func NewCancellableReader(ctx context.Context, reader *strings.Reader) *strings.Reader {
	return reader
}
