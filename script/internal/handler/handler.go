package handler

import (
	"project/pkg/coss"
	"project/pkg/logger"
	"project/pkg/wechat"
	"project/script/internal/service"
	"time"
)

type Config struct {
	Cos coss.CosConfig
	//Oss    coss.OssConfig
	Cdn    string
	Wechat struct {
		Appid  string
		Secret string
	}
	Robot struct {
		DingTalk   string
		WechatWork string
	}
}

type Handler struct {
	service *service.Service
	cos     coss.COS
	//oss         coss.OSS
	wechat      wechat.BasicAPI
	robotDing   string
	robotWechat string
}

func New(cfg *Config, srv *service.Service) *Handler {
	return &Handler{
		service: srv,
		cos:     coss.NewTCOS(&cfg.Cos),
		//oss:         coss.NewAliOSS(&cfg.Oss),
		wechat:      wechat.NewBasicAPI(cfg.Wechat.Appid, cfg.Wechat.Secret, logger.NewHttpClient(15*time.Second)),
		robotDing:   cfg.Robot.DingTalk,
		robotWechat: cfg.Robot.WechatWork,
	}
}
