package logger

import "context"

type __ctx__ struct {
	m map[string]string
	context.Context
}

func (c *__ctx__) Value(key any) any {
	if ks, ok := key.(string); ok {
		return c.m[ks]
	} else {
		return nil
	}
}

func NewCtxLog(tid, v1, v2, v3 string) (context.Context, Logger) {
	return &__ctx__{
			m: map[string]string{
				"trace_id": tid,
				"v1":       v1,
				"v2":       v2,
				"v3":       v3,
			},
			Context: context.Background(),
		}, &logger{
			TraceId: tid,
			V1:      v1,
			V2:      v2,
			V3:      v3,
		}
}

func FromContext(ctx context.Context) Logger {
	return &logger{
		TraceId: ctx.Value("trace_id"),
		V1:      ctx.Value("v1"),
		V2:      ctx.Value("v2"),
		V3:      ctx.Value("v3"),
	}
}
