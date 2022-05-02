package handler

import (
	"github.com/gin-gonic/gin"
	"project/pkg/web"
)

func register(r *gin.Engine, h *Handler) {
	r.GET("/", func(c *gin.Context) { c.String(OK, "OK") })
	r.Use(web.Cors, web.SetContext, web.Recover)

	r.POST("captcha", h.Captcha)
	{
		g := r.Group("user", web.AccessLog)
		g.POST("login", h.UserLogin)
		g.POST("logout", h.CheckAuth, h.UserLogout)
		g.POST("password", h.CheckAuth, h.UserPassword)
	}

	rbacNoLog := r.Group("", h.CheckAuth, h.RBAC)
	rbac := rbacNoLog.Group("", web.AccessLog)

	{
		g := rbac.Group("system")
		g.POST("action-menu/sync", h.SystemActionMenuSync(r))
		g.POST("action-menu/list", h.SystemActionMenuList)
		g.POST("action/update", h.SystemActionUpdate)
		g.POST("role/list", h.AdminRoleList)
		g.POST("role/option", h.AdminRoleOption)
		g.POST("role/create", h.AdminRoleCreate)
		g.POST("role/update", h.AdminRoleUpdate)
		g.POST("user/list", h.AdminUserList)
		g.POST("user/create", h.AdminUserCreate)
		g.POST("user/password", h.AdminUserPassword)
		g.POST("user/role", h.AdminUserRole)
		g.POST("user/status", h.AdminUserStatus)
	}

	{
		g := rbacNoLog.Group("upload")
		g.POST("image", h.UploadImage)
		g.POST("audio", h.UploadAudio)
		g.POST("video", h.UploadVideo)
		g.POST("pdf", h.UploadPDF)
	}

	{
		g := rbacNoLog.Group("crypto")
		g.POST("encrypt", h.CryptoEncrypt)
		g.POST("decrypt", h.CryptoDecrypt)
	}
}
