package handler

import (
	"github.com/gin-gonic/gin"
	"project/pkg/ginutil"
)

func register(r *gin.Engine, h *Handler) {
	r.GET("/", func(c *gin.Context) { c.String(OK, "OK") })
	r.Use(ginutil.SetContext, ginutil.Recover)

	api := r.Group("", ginutil.AccessLog)

	api.POST("wechat/login", h.WechatLogin)

	{
		g := api.Group("wechat", h.CheckAuth)
		g.POST("phone", h.WechatPhone)
		g.POST("user/profile", h.UserProfile)
		g.POST("user/info", h.UserInfo)
	}
}
