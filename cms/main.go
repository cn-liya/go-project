package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"project/cms/internal/handler"
	"project/cms/internal/service"
	"project/pkg/logger"
	"project/pkg/process"
	"time"
)

func setup() *http.Server {
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	var cfg struct {
		App struct {
			Env    string
			Logger logger.Config
		}
		Handler handler.Config
		Service service.Config
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal(err)
	}

	process.SetEnv(cfg.App.Env)
	logger.Setup(&cfg.App.Logger)
	gin.SetMode(gin.ReleaseMode)

	s := service.New(&cfg.Service)
	h := handler.Initialize(&cfg.Handler, s)
	server := &http.Server{
		Addr:    ":6000",
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
	l := logger.NewLogger(time.Now().Format(time.RFC3339Nano), process.GetIP(), process.GetEnv(), "")
	l.Info("start", nil, nil)
	process.Notify()
	err := server.Shutdown(context.Background())
	l.Info("stop", nil, err)
}
