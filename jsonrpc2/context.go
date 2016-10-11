package jsonrpc2

import "context"

type WithContext interface {
	Context() context.Context
	SetContext(ctx context.Context)
}

type Ctx struct {
	ctx context.Context
}

func (c *Ctx) Context() context.Context {
	return c.ctx
}

func (c *Ctx) SetContext(ctx context.Context) {
	c.ctx = ctx
}
