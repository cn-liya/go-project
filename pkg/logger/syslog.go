package logger

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

const (
	fileDir    = "docs/log/"
	latest     = fileDir + "app.log"
	fileFormat = fileDir + "20060102150405.log"
)

var (
	once       = &sync.Once{}
	appLog     = log.New(os.Stdout, "", 0)
	handle     = func(*columns) {}
	colorNum   int8
	lastHandle *os.File
)

func SetOutput(output string) {
	once.Do(func() {
		switch output {
		case "std":
			setLogToStdout()
		case "fmt":
			setLogToFormat()
		case "file":
			setLogToFile()
		}
	})
}

func setLogToStdout() {
	handle = func(c *columns) {
		enc := json.NewEncoder(appLog.Writer())
		enc.SetEscapeHTML(false)
		_ = enc.Encode(c)
	}
}

func setLogToFormat() {
	handle = func(c *columns) {
		b := bytes.NewBuffer(nil)
		enc := json.NewEncoder(b)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "\t")
		_ = enc.Encode(c)
		colorNum = (colorNum + 3) & 7 // 相邻日志使用不同颜色(黄青红蓝灰绿紫黑)
		appLog.Printf("\x1b[0;%dm%s\x1b[0m", colorNum+30, b)
	}
}

func setLogToFile() {
	file, err := os.OpenFile(latest, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	lastHandle = file
	appLog.SetOutput(file)
	handle = func(c *columns) {
		enc := json.NewEncoder(appLog.Writer())
		enc.SetEscapeHTML(false)
		_ = enc.Encode(c)
	}

	go func() {
		tick := time.Tick(7 * time.Second) //扫描频率
		for t := range tick {
			info, _ := lastHandle.Stat()
			if info.Size() > 500<<20 { //超过500M切割一次
				os.Rename(latest, t.Format(fileFormat)) //nolint （windows系统下无法重命名正在打开的文件）
				if f, e := os.OpenFile(latest, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); e == nil {
					appLog.SetOutput(f)
					lastHandle.Close() //nolint
					lastHandle = f
				}
			}
			dirs, _ := os.ReadDir(fileDir)
			if len(dirs) > 3 { //删除旧文件
				os.Remove(fileDir + dirs[0].Name()) //nolint
			}
		}
	}()
}
