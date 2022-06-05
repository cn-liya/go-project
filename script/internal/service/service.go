package service

import (
	"github.com/spf13/viper"
	"log"
	"project/pkg/cache"
	"project/pkg/db"
	"project/pkg/logger"
)

var config struct {
	App struct {
		IsProd bool
		Logger string
	}
	Mysql db.Mysql
	Redis cache.Redis
	Nsq   struct {
		Producer string
		Consumer string
	}
	Wechat struct {
		Appid  string
		Secret string
	}
	Robot struct {
		DingTalk   string
		WechatWork string
	}
}

func Initialize() {
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./script")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("viper.ReadInConfig error", err)
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("viper.Unmarshal error: ", err)
	}
	logger.SetOutput(config.App.Logger)
}
