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

var env, ip string

func init() {
	addrs, _ := net.InterfaceAddrs()
	for _, v := range addrs {
		if ipn, ok := v.(*net.IPNet); ok && !ipn.IP.IsLoopback() {
			if ipn.IP.To4() != nil {
				ip = ipn.IP.String()
				break
			}
		}
	}
}

func GetIP() string {
	return ip
}

func SetEnv(v string) {
	if v != EnvDev && v != EnvTest && v != EnvProd {
		log.Fatal("invalid env:", v)
	}
	env = v
}

func GetEnv() string {
	return env
}

// Notify 阻塞主进程，监听退出信息
func Notify() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
}
