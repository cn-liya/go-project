package logger

import (
	"bytes"
	"fmt"
	"os"
	"project/pkg/logger/internal/json"
)

var handle = func(*columns) {}

func setLogToFile(topic string) {
	writer := NewFileWriter(&Rotate{
		FileName:     "log/" + topic + "/app.log",
		MaxMegabytes: 100,
		MaxBackups:   2,
	})
	handle = func(c *columns) {
		buf := bytes.NewBuffer(nil)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		_ = enc.Encode(c)
		b := matchReplace(buf.Bytes())
		writer.Write(b) //nolint
	}
}

func setLogToStdout() {
	handle = func(c *columns) {
		buf := bytes.NewBuffer(nil)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		_ = enc.Encode(c)
		b := matchReplace(buf.Bytes())
		os.Stdout.Write(b) //nolint
	}
}

func setLogToFormat() {
	var colorNum int8
	handle = func(c *columns) {
		buf := bytes.NewBuffer(nil)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "\t")
		_ = enc.Encode(c)
		b := matchReplace(buf.Bytes())
		colorNum = (colorNum + 3) & 7 // 相邻日志使用不同颜色(黄青红蓝灰绿紫黑)
		fmt.Printf("\x1b[0;%dm%s\x1b[0m", colorNum+30, b)
	}
}
