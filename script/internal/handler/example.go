package handler

import (
	"project/pkg/core"
	"project/pkg/logger"
	"time"
)

func (h *Handler) Example1(ctx *core.Context) {
	l := logger.FromContext(ctx)
	count, err := h.service.ExampleLLen(ctx)
	if err != nil {
		l.Error("service.ExampleLLen error", nil, err)
		return
	}
	if count > 1000 {
		l.Warn("service.ExampleLLen long", nil, count)
	}
}

func (h *Handler) Example2(ctx *core.Context) {
	l := logger.FromContext(ctx)
	b, err := h.service.ExampleRPop(ctx)
	if err != nil {
		l.Error("service.ExampleRPop error", nil, err)
		time.Sleep(3 * time.Second)
		return
	}
	if len(b) == 0 {
		time.Sleep(2 * time.Second)
		return
	}
	// ... more logic
}
