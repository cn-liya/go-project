package handler

import (
	"project/pkg/core"
	"project/pkg/dingtalk"
	"project/pkg/logger"
	"project/pkg/wechatwork"
	"strconv"
	"time"
)

func (h *Handler) WechatServerToken(ctx *core.Context) {
	l := logger.FromContext(ctx)
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

func (h *Handler) ReportNewUser(ctx *core.Context) {
	l := logger.FromContext(ctx)
	now := time.Now()
	end := now.Format(time.DateOnly)
	begin := now.AddDate(0, 0, -1).Format(time.DateOnly)
	count, err := h.service.UserCount(ctx, begin, end)
	if err != nil {
		l.Error("service.UserCount error", nil, err)
		return
	}
	content := begin + "新增用户数：" + strconv.FormatInt(count, 10)
	l.Info("content", nil, content)
	_, _ = wechatwork.SendText(h.robotWechat, &wechatwork.Text{Content: content})
	_, _ = dingtalk.SendText(h.robotDing, &dingtalk.Text{Content: content}, nil)
}
