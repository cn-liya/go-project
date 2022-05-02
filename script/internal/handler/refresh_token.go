package handler

import (
	"project/pkg/logger"
	"project/pkg/random"
	"time"
)

func (h *Handler) WechatServerToken() {
	ctx, l := logger.New(random.UUID(), "RefreshToken", "WechatServerToken", "")
	ttl, err := h.service.WechatTokenTTL(ctx)
	if err != nil {
		l.Error("service.WechatTokenTTL error", nil, err)
		return
	}
	if ttl > 10*time.Minute {
		return
	}
	resp, err := h.wechat.GetAccessToken(ctx)
	if err != nil {
		l.Error("wechat.AccessToken error", nil, err)
		return
	}
	if resp.Errcode == 0 && resp.AccessToken != "" {
		err = h.service.WechatTokenSet(ctx, resp.AccessToken, time.Duration(resp.ExpiresIn)*time.Second)
		if err != nil {
			l.Error("service.WechatTokenSet error", nil, err)
		}
	} else {
		l.Warn("wechat.AccessToken fail", nil, resp)
	}
}
