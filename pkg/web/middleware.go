package web

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"project/pkg/logger"
	"project/pkg/random"
	"project/pkg/web/errcode"
	"runtime"
	"slices"
	"time"
)

func SetContext(c *gin.Context) {
	v0 := c.GetHeader("X-Request-Id")
	if v0 == "" {
		v0 = random.UUID()
	}
	c.Set("v0", v0)
	c.Set("v1", c.Request.URL.Path)
	c.Next()
}

func Recover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			buff := make([]byte, 2<<10)
			runtime.Stack(buff, false)
			logger.FromContext(c).Fatal("recover", err, bytes.TrimRight(buff, "\u0000"))
			c.AbortWithStatusJSON(errcode.InternalServerError.Response())
		}
	}()
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
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	w := &BodyLogWriter{
		ResponseWriter: c.Writer,
		body:           bytes.NewBuffer(nil),
	}
	c.Writer = w
	c.Next()
	logger.FromContext(c).Trace("access",
		gin.H{
			"query":     logger.SpreadMaps(c.Request.URL.Query()),
			"header":    logger.SpreadMaps(c.Request.Header),
			"body":      logger.Compress(body),
			"client_ip": c.ClientIP(),
		},
		gin.H{
			"body":   logger.Compress(w.body.Bytes()),
			"status": w.Status(),
		},
		begin)
}

/*
Cors 添加跨域头，或在nginx配置，两者取其一，不能同时存在。

	add_header Access-Control-Allow-Origin $http_origin always;
	add_header Access-Control-Allow-Methods 'GET, POST';
	add_header Access-Control-Allow-Headers 'Authorization, X-Request-Id';
	if ($request_method = 'OPTIONS') {
		return 204;
	}
*/
func Cors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
	c.Header("Access-Control-Allow-Methods", "GET, POST")
	c.Header("Access-Control-Allow-Headers", "Authorization, X-Request-Id")
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}

func NetworkLimit(ips []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if slices.Contains(ips, ip) {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(errcode.Forbidden.WithMsg("网络受限：" + ip))
	}
}
