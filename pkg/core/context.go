package core

import (
	"context"
	"os"
	"project/pkg/random"
	"sync"
)

type Context struct {
	m *sync.Map
	context.Context
}

func (c *Context) Value(key any) any {
	if v, ok := c.m.Load(key); ok {
		return v
	}
	return c.Context.Value(key)
}

func (c *Context) Set(key string, value any) {
	c.m.Store(key, value)
}

func (c *Context) Get(key string) (any, bool) {
	return c.m.Load(key)
}

func ContextWithCtx(ctx context.Context) *Context {
	return &Context{
		m:       &sync.Map{},
		Context: ctx,
	}
}

func ContextWithVal(v0, v1, v2, v3 string) *Context {
	if v0 == "" {
		v0 = random.UUID()
	}
	if v1 == "" && len(os.Args) > 1 {
		v1 = os.Args[1]
	}
	c := ContextWithCtx(context.Background())
	c.Set("v0", v0)
	c.Set("v1", v1)
	c.Set("v2", v2)
	c.Set("v3", v3)
	return c
}
