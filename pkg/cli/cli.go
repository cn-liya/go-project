package cli

import "project/pkg/core"

type HandlerFunc func(*core.Context)

func wrap(h HandlerFunc) func() {
	v2 := core.FuncName(h)
	return func() {
		ctx := core.ContextWithVal("", "", v2, "")
		defer core.Recover(ctx)
		h(ctx)
	}
}
