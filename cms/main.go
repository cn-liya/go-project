package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"project/cms/internal/handler"
	"project/cms/internal/service"
	"project/pkg/logger"
	"project/pkg/process"
)

func setup() *http.Server {
	var cfg struct {
		App struct {
			Env    string
			Logger logger.Config
		}
		Handler handler.Config
		Service service.Config
	}

	b, err := os.ReadFile("conf.yaml")
	if err != nil {
		log.Fatal(err)
	}
	if err = yaml.Unmarshal(b, &cfg); err != nil {
		log.Fatal(err)
	}

	process.SetEnv(cfg.App.Env)
	logger.Setup(&cfg.App.Logger)
	if process.GetEnv() != process.EnvDev {
		gin.SetMode(gin.ReleaseMode)
	}

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
	l := logger.New(process.GetHostname(), "cms", process.GetIP(), process.GetEnv())
	l.Info("start", nil, nil)
	process.Notify()
	err := server.Shutdown(context.Background())
	l.Info("stop", nil, err)
}
