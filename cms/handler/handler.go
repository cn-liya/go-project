package handler

import (
	"bytes"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"io"
	"net/http"
	"project/cms/internal/acl"
	"project/cms/internal/service"
	"project/pkg/logger"
	"project/pkg/util/captcha"
	"reflect"
	"runtime"
	"strings"
	"time"
)

type Config struct {
	Captcha string
	Cdn     string
}

type Handler struct {
	captcha string
	cdn     string
	service *service.Service
	drawer  *captcha.Drawer
}

func Initialize(cfg *Config, srv *service.Service) *gin.Engine {
	h := &Handler{
		captcha: cfg.Captcha,
		cdn:     cfg.Cdn,
		service: srv,
		drawer: captcha.NewDrawer("docs/fonts/AppleChancery.ttf",
			"docs/img/bg1.jpeg", ""),
	}
	r := gin.New()
	h.register(r)
	return r
}

//alias short for HttpStatusCode
const (
	OK                  = http.StatusOK                    //200: 成功
	InvalidParam        = http.StatusBadRequest            //400: 参数错误
	NotAuth             = http.StatusUnauthorized          //401: 登录失效
	Forbidden           = http.StatusForbidden             //403: 禁止操作
	NotFound            = http.StatusNotFound              //404: 目标不存在
	Conflict            = http.StatusConflict              //409: 数据已存在
	OverSize            = http.StatusRequestEntityTooLarge //413: 提交内容过大
	UnsupportedType     = http.StatusUnsupportedMediaType  //415: 错误的文件类型
	UnprocessableEntity = http.StatusUnprocessableEntity   //422: 数据格式错误或已过期
	Locked              = http.StatusLocked                //423: 资源被锁定
	RateLimit           = http.StatusTooManyRequests       //429: 请求频率限制
	ServerError         = http.StatusInternalServerError   //500: 服务端通用错误
	WrongResponse       = http.StatusBadGateway            //502: 响应错误
	ServiceUnavailable  = http.StatusServiceUnavailable    //503: 服务不可用
	GatewayTimeout      = http.StatusGatewayTimeout        //504: 请求错误
)

type RespErr struct {
	Msg    string `json:"msg"`
	Detail string `json:"detail,omitempty"`
}

var Empty = struct{}{}

func RespWithMsg(code int, msg string) (int, *RespErr) {
	return code, &RespErr{
		Msg: msg,
	}
}

func RespWrapErr(err error) (int, *RespErr) {
	code, msg, detail := ServerError, "系统繁忙", ""
	e := reflect.TypeOf(err).String()
	switch e {
	case "validator.ValidationErrors":
		code = InvalidParam
		msg = "参数错误"
		detail = err.Error()
	case "proto.RedisError":
		detail = "REDIS"
	case "nsq.ErrProtocol":
		detail = "NSQ"
	case "*errors.errorString":
		detail = "ERRORS"
	case "*url.Error":
		code = GatewayTimeout
		detail = "REQUEST"
	default:
		if strings.HasPrefix(e, "*json.") {
			code = WrongResponse
			detail = "RESPONSE"
		}
	}
	return code, &RespErr{Msg: msg, Detail: detail}
}

func Recover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			buff := make([]byte, 2<<10)
			runtime.Stack(buff, false)
			logger.FromContext(c).Fatal("recover", err, bytes.TrimRight(buff, "\u0000"))
			c.AbortWithStatusJSON(ServerError, &RespErr{
				Msg:    "系统繁忙",
				Detail: "RECOVER",
			})
		}
	}()
	c.Next()
}

func Cors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Trace-Id")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}

func SetContext(c *gin.Context) {
	tid := c.GetHeader("X-Trace-Id")
	if tid == "" {
		tid = base64.RawURLEncoding.EncodeToString(uuid.NewV4().Bytes())
	}
	c.Set("trace_id", tid)
	c.Set("v1", c.Request.Method+c.Request.URL.Path)
	c.Next()
}

type BodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *BodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func AccessLog(c *gin.Context) {
	begin := time.Now()
	body, _ := io.ReadAll(c.Request.Body)
	if len(body) > 0 {
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
	}

	w := &BodyLogWriter{
		ResponseWriter: c.Writer,
		body:           bytes.NewBuffer(nil),
	}
	c.Writer = w

	c.Next()

	logger.FromContext(c).Trace("access",
		gin.H{
			"query":     logger.SpreadMaps(c.Request.URL.Query()),
			"headers":   logger.SpreadMaps(c.Request.Header),
			"body":      logger.Compress(body),
			"client_ip": c.ClientIP(),
		},
		gin.H{
			"body":   logger.Compress(w.body.Bytes()),
			"status": w.Status(),
		}, begin)
}

func (h *Handler) AuthCheck(module string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(RespWithMsg(NotAuth, "Authorization Missing"))
			return
		}
		user, err := h.service.GetAdminToken(c, token)
		if err != nil {
			logger.FromContext(c).Error("service.GetUserToken error", token, err)
			c.AbortWithStatusJSON(RespWrapErr(err))
			return
		}
		if user.ID == 0 {
			c.AbortWithStatusJSON(RespWithMsg(NotAuth, "Authorization Expired"))
			return
		}
		if module != "" && user.Username != acl.SuperUser {
			rw := user.Authority[module]
			if rw == 0 || (rw == 1 && c.Request.Method != http.MethodGet) {
				c.AbortWithStatusJSON(RespWithMsg(Forbidden, ""))
				return
			}
		}
		c.Set("user", user)
		c.Set("v2", user.Username)
		c.Set("v3", c.ClientIP())
		c.Next()
	}
}

const (
	TimeFormat     = "2006-01-02 15:04:05"
	PureDateFormat = "20060102"
)
