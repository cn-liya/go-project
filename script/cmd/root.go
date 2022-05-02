package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"project/pkg/core"
	"project/pkg/db"
	"project/pkg/gredis"
	"project/pkg/logger"
	"project/pkg/process"
	"project/script/internal/handler"
	"strings"
)

var cfg struct {
	App struct {
		Env    string
		Logger logger.Config
	}
	Handler handler.Config
	Service struct {
		Mysql db.Config
		Redis gredis.Config
		Nsq   struct {
			Producer string
			Consumer string
		}
		//Rabbitmq rabbitmq.Config
	}
}

func init() {
	cobra.OnInitialize(func() {
		b, err := os.ReadFile("conf.yaml")
		if err != nil {
			log.Fatal(err)
		}
		if err = yaml.Unmarshal(b, &cfg); err != nil {
			log.Fatal(err)
		}
	})
	rootCmd.AddCommand(
		cronjobCmd,
		avatarToCdnCmd,
	)
}

var rootCmd = &cobra.Command{
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		process.SetEnv(cfg.App.Env)
		cfg.App.Logger.Topic = strings.ReplaceAll(cmd.Use, ":", "-")
		logger.Setup(&cfg.App.Logger)
		ctx := core.ContextWithVal(process.GetHostname(), cmd.Use, process.GetIP(), process.GetEnv())
		cmd.SetContext(ctx)
		logger.FromContext(ctx).Info("start", nil, nil)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		logger.FromContext(cmd.Context()).Info("stop", nil, nil)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
