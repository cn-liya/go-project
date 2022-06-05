package logger

import (
	"fmt"
	"time"
)

type Logger interface {
	Fatal(msg string, input, output any)
	Error(msg string, input, output any)
	Warn(msg string, input, output any)
	Info(msg string, input, output any)
	Trace(msg string, input, output any, begin time.Time)
}

type logger struct {
	TraceId any `json:"trace_id"`
	V1      any `json:"v1,omitempty"`
	V2      any `json:"v2,omitempty"`
	V3      any `json:"v3,omitempty"`
}

type level int8

type columns struct {
	*logger
	Level   level  `json:"level"`
	Time    string `json:"time"`
	Msg     string `json:"msg"`
	Input   any    `json:"input,omitempty"`
	Output  any    `json:"output,omitempty"`
	Elapsed int64  `json:"elapsed,omitempty"`
}

const (
	_ level = iota
	levelFatal
	levelError
	levelWarn
	levelInfo
)

func covert(val any) any {
	switch v := val.(type) {
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	case []byte:
		return string(v)
	default:
		return v
	}
}

func (l *logger) stash(level level, msg string, input, output any, et int64) {
	logs := &columns{
		logger:  l,
		Level:   level,
		Time:    time.Now().Format("2006/01/02-15:04:05.000000"),
		Msg:     msg,
		Input:   covert(input),
		Output:  covert(output),
		Elapsed: et,
	}
	handle(logs)
}

func (l *logger) Fatal(msg string, input, output any) {
	l.stash(levelFatal, msg, input, output, 0)
}

func (l *logger) Error(msg string, input, output any) {
	l.stash(levelError, msg, input, output, 0)
}

func (l *logger) Warn(msg string, input, output any) {
	l.stash(levelWarn, msg, input, output, 0)
}

func (l *logger) Info(msg string, input, output any) {
	l.stash(levelInfo, msg, input, output, 0)
}

func (l *logger) Trace(msg string, input, output any, begin time.Time) {
	l.stash(levelInfo, msg, input, output, time.Since(begin).Milliseconds())
}
