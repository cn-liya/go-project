package logger

// 参考"gopkg.in/natefinch/lumberjack.v2"包按需精简

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	megabyte          = 1 << 20
	defaultMaxSize    = 100 * megabyte
	defaultMaxBackups = 2
	defaultFilename   = "log/app.log"
	backupFormat      = "20060102150405.000"
)

type rotate struct {
	dir, prefix, suffix       string
	maxBackups, maxSize, size int
	file                      *os.File
	mu                        sync.Mutex
	async                     chan struct{}
}

type Rotate struct {
	FileName     string
	MaxMegabytes int
	MaxBackups   int
}

func NewFileWriter(cfg *Rotate) io.Writer {
	filename := cfg.FileName
	if filename == "" {
		filename = defaultFilename
	}
	dir := filepath.Dir(filename)
	_ = os.MkdirAll(dir, 0755)
	base := filepath.Base(filename)
	suffix := filepath.Ext(base)
	prefix := base[:len(base)-len(suffix)]
	l := &rotate{dir: dir, prefix: prefix, suffix: suffix}
	if err := l.open(); err != nil {
		log.Fatal(err)
	}
	stat, _ := l.file.Stat()
	l.size = int(stat.Size())
	if cfg.MaxMegabytes > 0 {
		l.maxSize = cfg.MaxMegabytes * megabyte
	} else {
		l.maxSize = defaultMaxSize
	}
	if cfg.MaxBackups > 0 {
		l.maxBackups = cfg.MaxBackups
	} else {
		l.maxBackups = defaultMaxBackups
	}
	l.async = make(chan struct{}, 1)
	go l.backup()
	return l
}

func (l *rotate) Write(b []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.size > l.maxSize {
		l.rotate()
	}
	l.size += len(b)
	return l.file.Write(b)
}

// nolint
func (l *rotate) rotate() {
	l.file.Close()
	os.Rename(l.file.Name(), l.dir+"/"+l.prefix+"."+time.Now().Format(backupFormat)+l.suffix)
	l.open()
	select {
	case l.async <- struct{}{}:
	default:
	}
}

func (l *rotate) open() error {
	file, err := os.OpenFile(l.dir+"/"+l.prefix+l.suffix, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	l.file = file
	l.size = 0
	return nil
}

func (l *rotate) backup() {
	for range l.async {
		files, _ := os.ReadDir(l.dir)
		count := len(files)
		if count <= l.maxBackups+1 {
			continue
		}
		backups := make([]string, 0, count)
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			if name := f.Name(); l.isBackupFile(name) {
				backups = append(backups, name)
			}
		}
		del := len(backups) - l.maxBackups
		for i := 0; i < del; i++ {
			os.Remove(l.dir + "/" + backups[i]) //nolint
		}
	}
}

func (l *rotate) isBackupFile(name string) bool {
	if !strings.HasPrefix(name, l.prefix) || !strings.HasSuffix(name, l.suffix) ||
		name == filepath.Base(l.file.Name()) {
		return false
	}
	_, err := time.Parse(backupFormat, name[len(l.prefix)+1:len(name)-len(l.suffix)])
	return err == nil
}
