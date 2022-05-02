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

const (
	_ int8 = iota
	levelFatal
	levelError
	levelWarn
	levelInfo
)
const timeFormat = "2006/01/02-15:04:05.000000"

type logger struct {
	tid string
	v1  string
	v2  string
	v3  string
}

type columns struct {
	TraceId string `json:"trace_id"`
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
		TraceId: l.tid,
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
	l.stash(levelInfo, msg, input, output, time.Since(begin).Milliseconds())
}
