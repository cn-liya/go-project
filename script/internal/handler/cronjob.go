package handler

import (
	"project/pkg/dingtalk"
	"project/pkg/logger"
	"project/pkg/random"
	"project/pkg/wechatwork"
	"strconv"
	"time"
)

func (h *Handler) ReportNewUser() {
	now := time.Now()
	end := now.Format(time.DateOnly)
	begin := now.AddDate(0, 0, -1).Format(time.DateOnly)
	ctx, l := logger.New(random.UUID(), "Cronjob", "ReportNewUser", begin)
	count, err := h.service.UserCount(ctx, begin, end)
	if err != nil {
		l.Error("service.UserCount error", nil, err)
		return
	}
	content := begin + "新增用户数：" + strconv.FormatInt(count, 10)
	_, _ = wechatwork.SendText(h.robotWechat, &wechatwork.Text{Content: content})
	_, _ = dingtalk.SendText(h.robotDing, &dingtalk.Text{Content: content}, nil)
}
