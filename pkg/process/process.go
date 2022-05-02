package process

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	EnvDev  = "development"
	EnvTest = "testing"
	EnvProd = "production"
)

var env, ip, hostname string

func init() {
	if v := os.Getenv("APP_ENV"); v != "" {
		SetEnv(v)
	}
	addrs, _ := net.InterfaceAddrs()
	for _, v := range addrs {
		if ipn, ok := v.(*net.IPNet); ok && !ipn.IP.IsLoopback() {
			if ipn.IP.To4() != nil {
				ip = ipn.IP.String()
				break
			}
		}
	}
	hostname = os.Getenv("HOSTNAME")
}

func SetEnv(v string) {
	if v != EnvDev && v != EnvTest && v != EnvProd {
		log.Fatal("invalid env:", v)
	}
	env = v
}

func GetEnv() string { return env }

func GetIP() string { return ip }

func GetHostname() string { return hostname }

// Notify 阻塞主进程，直到捕获退出通知信号
func Notify() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
}
