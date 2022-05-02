package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"project/api/internal/service"
	"project/model/cache"
	"project/pkg/ades"
	"project/pkg/logger"
	"project/pkg/process"
	"project/pkg/web"
	"project/pkg/web/errcode"
	"project/pkg/wechat"
	"strings"
	"time"
)

type Config struct {
	Ips    []string
	Cdn    string
	Wechat struct {
		Appid  string
		Secret string
	}
	Aes struct {
		Key string
		IV  string
	}
}

type Handler struct {
	service *service.Service
	wechat  wechat.FullAPI
	cdn     string
	aes     *ades.Cipher
}

func Initialize(cfg *Config, srv *service.Service) *gin.Engine {
	aes, err := ades.NewAesCipher([]byte(cfg.Aes.Key), []byte(cfg.Aes.IV))
	if err != nil {
		log.Fatal(err)
	}
	client := logger.NewHttpClient(8 * time.Second)
	h := &Handler{
		service: srv,
		wechat:  wechat.NewFullAPI(cfg.Wechat.Appid, cfg.Wechat.Secret, client, srv.WechatToken),
		cdn:     cfg.Cdn,
		aes:     aes,
	}
	r := gin.New()
	if process.GetEnv() != process.EnvProd { // 非正式环境全局注册IP白名单中间件
		r.Use(web.NetworkLimit(cfg.Ips))
	}
	register(r, h)
	return r
}

const OK = http.StatusOK

var Empty = struct{}{}

func (h *Handler) CheckAuth(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(errcode.Unauthorized.Response())
		return
	}
	user, err := h.service.UserTokenGet(c, token)
	if err != nil {
		logger.FromContext(c).Error("service.UserTokenGet error", token, err)
		c.AbortWithStatusJSON(errcode.FromError(err))
		return
	}
	if user.ID == 0 {
		c.AbortWithStatusJSON(errcode.Unauthorized.Response())
		return
	}
	c.Set("user", user)
	c.Set("v2", user.Openid)
	c.Set("v3", user.Unionid)
	c.Next()
}

func (h *Handler) CdnLink(path string) string {
	if path == "" || strings.HasPrefix(path, "http") {
		return path
	}
	return h.cdn + path
}

func getAuth(c *gin.Context) *cache.UserToken {
	v, _ := c.Get("user")
	return v.(*cache.UserToken)
}
