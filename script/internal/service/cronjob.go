package service

import (
	"context"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"project/model"
	"project/pkg/cache"
	"project/pkg/db"
	"project/pkg/dingtalk"
	"project/pkg/logger"
	"project/pkg/wechat"
	"project/pkg/wechatwork"
	"strconv"
	"strings"
	"time"
)

type Cronjob struct {
	mysql       *gorm.DB
	redis       *redis.Client
	wechat      wechat.ServerAPI
	robotDing   string
	robotWechat string
}

func NewCronJob() *Cronjob {
	s := &Cronjob{
		mysql:       db.NewMysqlDB(&config.Mysql),
		redis:       cache.NewRedisClient(&config.Redis),
		robotDing:   config.Robot.DingTalk,
		robotWechat: config.Robot.WechatWork,
	}
	s.wechat = wechat.NewServerAPI(logger.NewHttpClient(30*time.Second), func(ctx context.Context) (string, error) {
		return s.redis.Get(ctx, model.KeyWechatToken).Result()
	})
	return s
}

func (s *Cronjob) GetDailySummary() {
	ctx, l := logger.NewCtxLog(uuid.NewV4().String(), "Cronjob", "GetDailySummary", "")
	yesterday := time.Now().AddDate(0, 0, -1).Format(wechat.DateFormat)
	args := &wechat.DatacubeArgs{
		BeginDate: yesterday,
		EndDate:   yesterday,
	}
	resp, err := s.wechat.GetDailySummary(ctx, args)
	if err != nil {
		resp, err = s.wechat.GetDailySummary(ctx, args)
	}
	if err != nil {
		l.Error("wechat.GetDailySummary error", args, err)
		return
	}
	if resp.Errcode != 0 || len(resp.List) == 0 {
		l.Warn("wechat.GetDailySummary fail", args, resp)
		return
	}
	data := &model.AnalysisDailySummary{
		RefDate:    resp.List[0].RefDate,
		VisitTotal: resp.List[0].VisitTotal,
		SharePv:    resp.List[0].ShareUv,
		ShareUv:    resp.List[0].SharePv,
	}
	err = s.mysql.WithContext(ctx).Create(data).Error
	if err != nil {
		l.Error("mysql.Create error", data, err)
	}
	_, _ = dingtalk.SendText(s.robotDing, &dingtalk.Text{
		Content: "小程序累计访问人数已达：" + strconv.Itoa(data.VisitTotal),
	}, nil)
}

func (s *Cronjob) GetDailyVisitTrend() {
	ctx, l := logger.NewCtxLog(uuid.NewV4().String(), "Cronjob", "GetDailyVisitTrend", "")
	yesterday := time.Now().AddDate(0, 0, -1).Format(wechat.DateFormat)
	args := &wechat.DatacubeArgs{
		BeginDate: yesterday,
		EndDate:   yesterday,
	}
	resp, err := s.wechat.GetDailyVisitTrend(ctx, args)
	if err != nil {
		resp, err = s.wechat.GetDailyVisitTrend(ctx, args)
	}
	if err != nil {
		l.Error("wechat.GetDailyVisitTrend error", args, err)
		return
	}
	if resp.Errcode != 0 || len(resp.List) == 0 {
		l.Warn("wechat.GetDailyVisitTrend fail", args, resp)
		return
	}
	data := &model.AnalysisDailyTrend{
		RefDate:         resp.List[0].RefDate,
		SessionCnt:      resp.List[0].SessionCnt,
		VisitPv:         resp.List[0].VisitPv,
		VisitUv:         resp.List[0].VisitUv,
		VisitUvNew:      resp.List[0].VisitUvNew,
		StayTimeUv:      resp.List[0].StayTimeUv,
		StayTimeSession: resp.List[0].StayTimeSession,
		VisitDepth:      resp.List[0].VisitDepth,
	}
	err = s.mysql.WithContext(ctx).Create(data).Error
	if err != nil {
		l.Error("mysql.Create error", data, err)
	}
	text := &strings.Builder{}
	text.WriteString(data.RefDate)
	text.WriteString("\n访问PV:")
	text.WriteString(strconv.Itoa(data.VisitPv))
	text.WriteString("\n访问UV:")
	text.WriteString(strconv.Itoa(data.VisitUv))
	text.WriteString("\n新增用户:")
	text.WriteString(strconv.Itoa(data.VisitUvNew))
	_, _ = wechatwork.SendText(s.robotWechat, &wechatwork.Text{
		Content: text.String(),
	})
}

func (s *Cronjob) GetWeeklyVisitTrend() {
	ctx, l := logger.NewCtxLog(uuid.NewV4().String(), "Cronjob", "GetWeeklyVisitTrend", "")
	args := &wechat.DatacubeArgs{
		BeginDate: time.Now().AddDate(0, 0, -7).Format(wechat.DateFormat),
		EndDate:   time.Now().AddDate(0, 0, -1).Format(wechat.DateFormat),
	}
	resp, err := s.wechat.GetWeeklyVisitTrend(ctx, args)
	if err != nil {
		resp, err = s.wechat.GetWeeklyVisitTrend(ctx, args)
	}
	if err != nil {
		l.Error("wechat.GetWeeklyVisitTrend error", args, err)
		return
	}
	if resp.Errcode != 0 || len(resp.List) == 0 {
		l.Warn("wechat.GetWeeklyVisitTrend fail", args, resp)
		return
	}
	data := &model.AnalysisWeeklyTrend{
		RefDate:         resp.List[0].RefDate,
		SessionCnt:      resp.List[0].SessionCnt,
		VisitPv:         resp.List[0].VisitPv,
		VisitUv:         resp.List[0].VisitUv,
		VisitUvNew:      resp.List[0].VisitUvNew,
		StayTimeUv:      resp.List[0].StayTimeUv,
		StayTimeSession: resp.List[0].StayTimeSession,
		VisitDepth:      resp.List[0].VisitDepth,
	}
	err = s.mysql.WithContext(ctx).Create(data).Error
	if err != nil {
		l.Error("mysql.Create error", data, err)
	}
}
