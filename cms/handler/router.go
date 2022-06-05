package handler

import (
	"github.com/gin-gonic/gin"
	"project/cms/internal/acl"
)

func (h *Handler) register(r *gin.Engine) {
	r.GET("ping", func(c *gin.Context) {
		c.String(OK, "pong")
	})
	r.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(NotFound)
	})
	r.Use(Recover) // ,Cors (if nginx don't add cors header)

	if gin.Mode() != gin.ReleaseMode {
		r.GET("modules", func(c *gin.Context) {
			c.JSON(OK, gin.H{"list": acl.Modules})
		})
	}
	r.GET("captcha", h.Captcha)

	api := r.Group("", SetContext, AccessLog)
	api.POST("user/login", h.UserLogin)
	api.DELETE("user/logout", h.AuthCheck(""), h.UserLogout)
	api.PUT("user/password", h.AuthCheck(""), h.ModifyUserPassword)

	{
		admin := api.Group("admin", h.AuthCheck(acl.ModuleAdmin))
		admin.GET("list", h.AdminList)
		admin.POST("account", h.CreateAdmin)
		admin.PUT("password", h.ResetAdminPassword)
		admin.PUT("authority", h.AssignAdminAuthority)
		admin.PUT("status", h.SwitchAdminStatus)
	}

	{
		analysis := api.Group("analysis", h.AuthCheck(acl.ModuleAnalysis))
		analysis.GET("daily-summary", h.AnalysisDailySummary)
		analysis.GET("daily-trend", h.AnalysisDailyTrend)
		analysis.GET("weekly-trend", h.AnalysisWeeklyTrend)
	}
}
