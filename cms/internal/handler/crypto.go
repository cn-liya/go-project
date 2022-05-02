package handler

import (
	"github.com/gin-gonic/gin"
	"project/cms/internal/proto"
	"project/pkg/web/errcode"
)

func (h *Handler) CryptoEncrypt(c *gin.Context) {
	var r proto.CryptoArgs
	_ = c.ShouldBindJSON(&r)
	if l := len([]rune(r.Text)); l < 2 || l > 100 {
		c.JSON(errcode.InvalidParam.WithMsg("请输入2~100字符明文"))
		return
	}
	c.JSON(OK, &proto.CryptoResp{Result: h.aes.EncryptCTR([]byte(r.Text))})
}

func (h *Handler) CryptoDecrypt(c *gin.Context) {
	var r proto.CryptoArgs
	_ = c.ShouldBindJSON(&r)
	if l := len(r.Text); l < 4 || l > 400 {
		c.JSON(errcode.InvalidParam.WithMsg("请输入4~400字符密文"))
		return
	}
	result, _ := h.aes.DecryptCTR(r.Text)
	c.JSON(OK, &proto.CryptoResp{Result: string(result)})
}
