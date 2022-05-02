package core

import (
	"bytes"
	"fmt"
	"project/pkg/logger"
	"reflect"
	"runtime"
	"strings"
)

func FuncName(f any) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	if i := strings.LastIndex(name, "."); i > -1 {
		name = name[i+1:]
	}
	if i := strings.Index(name, "-"); i > 0 {
		name = name[:i]
	}
	return name
}

func Recover(ctx *Context) {
	if err := recover(); err != nil {
		buf := make([]byte, 2<<10)
		runtime.Stack(buf, false)
		logger.FromContext(ctx).Fatal("recover", err, bytes.TrimRight(buf, "\u0000"))
	}
}

func RecoverE(ctx *Context, e *error) {
	if err := recover(); err != nil {
		buf := make([]byte, 2<<10)
		runtime.Stack(buf, false)
		logger.FromContext(ctx).Fatal("recover", err, bytes.TrimRight(buf, "\u0000"))
		*e = fmt.Errorf("PANIC! %v", err)
	}
}
