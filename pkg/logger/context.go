package logger

import "context"

type __ctx__ struct {
	m map[string]string
	context.Context
}

func (c *__ctx__) Value(key any) any {
	if ks, ok := key.(string); ok {
		if v, exist := c.m[ks]; exist {
			return v
		}
	}
	return c.Context.Value(key)
}

func FromContext(ctx context.Context) Logger {
	l := &logger{}
	l.tid, _ = ctx.Value("trace_id").(string)
	l.v1, _ = ctx.Value("v1").(string)
	l.v2, _ = ctx.Value("v2").(string)
	l.v3, _ = ctx.Value("v3").(string)
	return l
}

func NewContext(tid, v1, v2, v3 string) context.Context {
	return &__ctx__{
		map[string]string{
			"trace_id": tid,
			"v1":       v1,
			"v2":       v2,
			"v3":       v3,
		},
		context.Background(),
	}
}

func NewLogger(tid, v1, v2, v3 string) Logger {
	return &logger{
		tid: tid,
		v1:  v1,
		v2:  v2,
		v3:  v3,
	}
}

func New(tid, v1, v2, v3 string) (context.Context, Logger) {
	return NewContext(tid, v1, v2, v3), NewLogger(tid, v1, v2, v3)
}
