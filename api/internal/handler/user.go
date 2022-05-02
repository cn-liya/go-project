package handler

import (
	"github.com/gin-gonic/gin"
	"project/api/internal/proto"
	"project/model/cache"
	"project/model/entity"
	"project/model/queue"
	"project/pkg/logger"
	"project/pkg/web/errcode"
)

func (h *Handler) WechatLogin(c *gin.Context) {
	var r proto.CodeArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(errcode.InvalidParam.Response())
		return
	}
	resp, err := h.wechat.Code2Session(c, r.Code)
	if err != nil {
		logger.FromContext(c).Error("wechat.Code2Session error", r.Code, err)
		c.JSON(errcode.FromError(err))
		return
	}
	if resp.Openid == "" {
		logger.FromContext(c).Warn("wechat.Code2Session fail", r.Code, resp)
		c.JSON(errcode.CodeExpiredOrWrong.Response())
		return
	}
	c.Set("v2", resp.Openid)
	c.Set("v3", resp.Unionid)
	uid, err := h.service.UserSave(c, &entity.User{
		Openid:  resp.Openid,
		Unionid: resp.Unionid,
	})
	if err != nil {
		logger.FromContext(c).Error("service.UserSave error", nil, err)
		c.JSON(errcode.FromError(err))
		return
	}
	token, err := h.service.UserTokenSet(c, &cache.UserToken{
		ID:         uid,
		Openid:     resp.Openid,
		Unionid:    resp.Unionid,
		SessionKey: resp.SessionKey,
	})
	if err != nil {
		logger.FromContext(c).Error("service.UserTokenSet error", uid, err)
		c.JSON(errcode.FromError(err))
		return
	}
	c.JSON(OK, &proto.LoginResp{
		Token:   token,
		Openid:  resp.Openid,
		Unionid: resp.Unionid,
	})
}

func (h *Handler) WechatPhone(c *gin.Context) {
	var r proto.CodeArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(errcode.InvalidParam.Response())
		return
	}
	resp, err := h.wechat.GetUserPhoneNumber(c, r.Code)
	if err != nil {
		logger.FromContext(c).Error("wechat.GetUserPhoneNumber error", r.Code, err)
		c.JSON(errcode.FromError(err))
		return
	}
	if resp.PhoneInfo == nil || resp.PhoneInfo.PhoneNumber == "" {
		logger.FromContext(c).Warn("wechat.GetUserPhoneNumber fail", r.Code, resp)
		c.JSON(errcode.CodeExpiredOrWrong.Response())
		return
	}
	user := getAuth(c)
	err = h.service.UserUpdate(c, &entity.User{
		ID:          user.ID,
		PhoneNumber: h.aes.EncryptCTR([]byte(resp.PhoneInfo.PhoneNumber)),
	})
	if err != nil {
		logger.FromContext(c).Error("service.SaveUserPhone error", resp.PhoneInfo, err)
		c.JSON(errcode.FromError(err))
		return
	}
	c.JSON(OK, &proto.WechatPhoneResp{PhoneNumber: resp.PhoneInfo.PhoneNumber})
}

func (h *Handler) UserProfile(c *gin.Context) {
	var r proto.UserProfileArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(errcode.InvalidParam.Response())
		return
	}
	user := getAuth(c)
	err := h.service.UserUpdate(c, &entity.User{
		ID:        user.ID,
		Nickname:  r.Nickname,
		AvatarURL: r.AvatarURL,
	})
	if err != nil {
		logger.FromContext(c).Error("service.UserUpdate error", &r, err)
		c.JSON(errcode.FromError(err))
		return
	}
	if err = h.service.AvatarToCdnAsync(&queue.AvatarToCdn{
		ID:        user.ID,
		Openid:    user.Openid,
		AvatarURL: r.AvatarURL,
	}); err != nil {
		logger.FromContext(c).Error("service.AvatarToCdnAsync error", nil, err)
	}
	c.JSON(OK, Empty)
}

func (h *Handler) UserInfo(c *gin.Context) {
	user := getAuth(c)
	info, err := h.service.UserFindByID(c, user.ID)
	if err != nil {
		logger.FromContext(c).Error("service.UserFindByID error", user.ID, err)
		c.JSON(errcode.FromError(err))
		return
	}
	resp := &proto.UserInfoResp{
		Nickname:  info.Nickname,
		AvatarURL: h.CdnLink(info.AvatarURL),
	}
	if info.PhoneNumber != "" {
		phone, _ := h.aes.DecryptCTR(info.PhoneNumber)
		resp.PhoneNumber = string(phone)
	}
	c.JSON(OK, resp)
}
