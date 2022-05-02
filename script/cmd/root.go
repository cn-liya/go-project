package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"project/pkg/db"
	"project/pkg/gredis"
	"project/pkg/logger"
	"project/pkg/process"
	"project/script/internal/handler"
	"strings"
	"time"
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
	}
}

func init() {
	cobra.OnInitialize(func() {
		viper.SetConfigName("conf")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal(err)
		}
		if err := viper.Unmarshal(&cfg); err != nil {
			log.Fatal(err)
		}
	})
	rootCmd.AddCommand(
		cronjobCmd,
		refreshTokenCmd,
		avatar2cdnCmd,
	)
}

var rootCmd = &cobra.Command{
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		process.SetEnv(cfg.App.Env)
		cfg.App.Logger.Topic = strings.ReplaceAll(cmd.Use, ":", "-")
		logger.Setup(&cfg.App.Logger)
		ctx, l := logger.New(time.Now().Format(time.RFC3339Nano), process.GetIP(), process.GetEnv(), cmd.Use)
		cmd.SetContext(ctx)
		l.Info("start", nil, nil)
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
