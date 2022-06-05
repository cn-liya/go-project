package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"project/api/handler"
	"project/api/internal/service"
	"project/pkg/logger"
	"syscall"
	"time"
)

func setup() *http.Server {
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./api")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("viper.ReadInConfig error", err)
	}

	var cfg struct {
		App struct {
			Mode   string
			Logger string
		}
		Handler handler.Config
		Service service.Config
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("viper.Unmarshal error: ", err)
	}

	gin.SetMode(cfg.App.Mode)
	logger.SetOutput(cfg.App.Logger)
	rand.Seed(time.Now().UnixNano())

	s := service.New(&cfg.Service)
	h := handler.Initialize(&cfg.Handler, s)
	server := &http.Server{
		Addr:    ":8000",
		Handler: h,
	}
	return server
}

func main() {
	server := setup()
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	c := <-quit
	log.Println("signal.Notify: ", c.String())

	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel() // 带超时控制，等待所有协程退出，或10秒强制退出
	ctx := context.Background() // 不带超时控制，等待所有协程退出
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}
	log.Println("Server Exit...")
}
