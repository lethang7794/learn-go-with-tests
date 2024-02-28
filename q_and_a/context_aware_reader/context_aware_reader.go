package context_aware_reader

import (
	"context"
	"io"
)

type readerCtx struct {
	ctx      context.Context
	delegate io.Reader
}

func (c readerCtx) Read(p []byte) (n int, err error) {
	if ctxErr := c.ctx.Err(); ctxErr != nil {
		return 0, ctxErr
	}
	return c.delegate.Read(p)
}

func NewCancellableReader(ctx context.Context, reader io.Reader) io.Reader {
	return &readerCtx{
		ctx:      ctx,
		delegate: reader,
	}
}
