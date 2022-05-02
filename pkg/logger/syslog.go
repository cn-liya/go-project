package logger

import (
	"context"
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
	v0 string
	v1 string
	v2 string
	v3 string
}

type columns struct {
	V0      string `json:"v0"`
	V1      string `json:"v1"`
	V2      string `json:"v2"`
	V3      string `json:"v3"`
	Level   int8   `json:"level"`
	Time    string `json:"time"`
	Msg     string `json:"msg"`
	Input   any    `json:"input,omitempty"`
	Output  any    `json:"output,omitempty"`
	Elapsed int64  `json:"elapsed,omitempty"`
}

const (
	timeFormat = "2006/01/02-15:04:05.000000"

	_ int8 = iota
	levelFatal
	levelError
	levelWarn
	levelInfo
	levelTrace
)

func convert(val any) any {
	switch v := val.(type) {
	case error:
		return v.Error()
	case fmt.Stringer:
		return Compress(v.String())
	case string:
		return Compress(v)
	case []byte:
		return Compress(v)
	default:
		return v
	}
}

func (l *logger) stash(level int8, msg string, input, output any, et int64) {
	c := &columns{
		V0:      l.v0,
		V1:      l.v1,
		V2:      l.v2,
		V3:      l.v3,
		Level:   level,
		Time:    time.Now().Format(timeFormat),
		Msg:     msg,
		Input:   convert(input),
		Output:  convert(output),
		Elapsed: et,
	}
	handle(c)
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
	l.stash(levelTrace, msg, input, output, time.Since(begin).Milliseconds())
}

func New(v0, v1, v2, v3 string) Logger {
	return &logger{v0, v1, v2, v3}
}

func FromContext(c context.Context) Logger {
	l := &logger{}
	l.v0, _ = c.Value("v0").(string)
	l.v1, _ = c.Value("v1").(string)
	l.v2, _ = c.Value("v2").(string)
	l.v3, _ = c.Value("v3").(string)
	return l
}
