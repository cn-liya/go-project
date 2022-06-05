package handler

import (
	"github.com/gin-gonic/gin"
	"project/api/internal/proto"
	"project/model"
	"project/pkg/logger"
)

type LoginArgs struct {
	JsCode string `json:"js_code" binding:"required"`
}
type LoginResp struct {
	Token   string `json:"token"`
	Openid  string `json:"openid"`
	Unionid string `json:"unionid"`
}

func (h *Handler) WechatLogin(c *gin.Context) {
	var r LoginArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(RespWithMsg(InvalidParam, err.Error()))
		return
	}

	user, err := h.service.WechatLogin(c, r.JsCode)
	if err != nil {
		logger.FromContext(c).Error("service.WechatLogin error", r.JsCode, err)
		c.JSON(RespWrapErr(err))
		return
	}

	if user.Openid == "" {
		c.JSON(RespWithMsg(UnprocessableEntity, "Invalid Or Expired"))
		return
	}
	token, err := h.service.SetUserToken(c, user)
	if err != nil {
		logger.FromContext(c).Error("service.SetUserToken error", user, err)
		c.JSON(RespWrapErr(err))
		return
	}

	c.JSON(OK, &LoginResp{
		Token:   token,
		Openid:  user.Openid,
		Unionid: user.Unionid,
	})
}

type WechatPhoneArgs struct {
	Code string `json:"code" binding:"required"`
}
type WechatPhoneResp struct {
	PhoneNumber string `json:"phone_number"`
}

func (h *Handler) WechatPhone(c *gin.Context) {
	var r WechatPhoneArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(RespWithMsg(InvalidParam, err.Error()))
		return
	}
	phone, err := h.service.WechatPhone(c, r.Code)
	if err != nil {
		logger.FromContext(c).Error("service.WechatPhone error", r.Code, err)
		c.JSON(RespWrapErr(err))
		return
	}
	if phone == "" {
		c.JSON(RespWithMsg(UnprocessableEntity, "Invalid Or Expired"))
		return
	}
	u, _ := c.Get("user")
	user := u.(*proto.UserToken)
	err = h.service.UpdateUser(c, &model.User{
		ID:          user.ID,
		PhoneNumber: phone,
	})
	if err != nil {
		logger.FromContext(c).Error("service.SaveUserPhone error", phone, err)
		c.JSON(RespWrapErr(err))
		return
	}
	c.JSON(OK, &WechatPhoneResp{PhoneNumber: phone})
}

type UserInfoArgs struct {
	Nickname  string `json:"nickname" binding:"required"`
	AvatarURL string `json:"avatar_url" binding:"required"`
}

func (h *Handler) SaveUserInfo(c *gin.Context) {
	var r UserInfoArgs
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(RespWithMsg(InvalidParam, err.Error()))
		return
	}
	u, _ := c.Get("user")
	user := u.(*proto.UserToken)
	err := h.service.UpdateUser(c, &model.User{
		ID:        user.ID,
		Nickname:  r.Nickname,
		AvatarURL: r.AvatarURL,
	})
	if err != nil {
		logger.FromContext(c).Error("service.SaveUserInfo error", &r, err)
		c.JSON(RespWrapErr(err))
		return
	}
	c.JSON(OK, Empty)
}

type GetUserInfoResp struct {
	PhoneNumber string `json:"phone_number"`
	Nickname    string `json:"nickname"`
	AvatarURL   string `json:"avatar_url"`
}

func (h *Handler) GetUserInfo(c *gin.Context) {
	u, _ := c.Get("user")
	user := u.(*proto.UserToken)
	info, err := h.service.QueryUser(c, user.ID)
	if err != nil {
		logger.FromContext(c).Error("service.QueryUser error", user.ID, err)
		c.JSON(RespWrapErr(err))
		return
	}
	c.JSON(OK, &GetUserInfoResp{
		PhoneNumber: info.PhoneNumber,
		Nickname:    info.Nickname,
		AvatarURL:   info.AvatarURL,
	})
}
