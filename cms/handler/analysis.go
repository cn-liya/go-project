package handler

import (
	"github.com/gin-gonic/gin"
	"project/cms/internal/proto"
	"project/pkg/logger"
	"time"
)

type AnalysisArgs struct {
	Begin string `form:"begin" binding:"required,datetime=20060102"`
	End   string `form:"end" binding:"required,datetime=20060102"`
}

type AnalysisResp struct {
	List any `json:"list"`
}

func ValidateDateRange(begin, end string, max int64) bool {
	t1, _ := time.Parse(PureDateFormat, begin)
	t2, _ := time.Parse(PureDateFormat, end)
	return (t2.Unix()-t1.Unix())/86400 < max
}

func (h *Handler) AnalysisDailySummary(c *gin.Context) {
	var r AnalysisArgs
	if err := c.ShouldBindQuery(&r); err != nil || r.Begin > r.End {
		c.JSON(RespWithMsg(InvalidParam, "请选择正确的时间范围:yyyymmdd"))
		return
	}
	if !ValidateDateRange(r.Begin, r.End, 90) {
		c.JSON(RespWithMsg(InvalidParam, "仅支持查询90天以内的数据"))
		return
	}
	list, err := h.service.AnalysisDailySummary(c, &proto.DateRange{
		Begin: r.Begin,
		End:   r.End,
	})
	if err != nil {
		logger.FromContext(c).Error("service.AnalysisDailySummary error", &r, err)
		c.JSON(RespWrapErr(err))
		return
	}
	c.JSON(OK, &AnalysisResp{
		List: list,
	})
}

func (h *Handler) AnalysisDailyTrend(c *gin.Context) {
	var r AnalysisArgs
	if err := c.ShouldBindQuery(&r); err != nil || r.Begin > r.End {
		c.JSON(RespWithMsg(InvalidParam, "请选择正确的时间范围:yyyymmdd"))
		return
	}
	if !ValidateDateRange(r.Begin, r.End, 90) {
		c.JSON(RespWithMsg(InvalidParam, "仅支持查询90天以内的数据"))
		return
	}
	list, err := h.service.AnalysisDailyTrend(c, &proto.DateRange{
		Begin: r.Begin,
		End:   r.End,
	})
	if err != nil {
		logger.FromContext(c).Error("service.AnalysisDailyTrend error", &r, err)
		c.JSON(RespWrapErr(err))
		return
	}
	c.JSON(OK, &AnalysisResp{
		List: list,
	})
}

func (h *Handler) AnalysisWeeklyTrend(c *gin.Context) {
	var r AnalysisArgs
	if err := c.ShouldBindQuery(&r); err != nil || r.Begin > r.End {
		c.JSON(RespWithMsg(InvalidParam, "请选择正确的时间范围:yyyymmdd"))
		return
	}
	if !ValidateDateRange(r.Begin, r.End, 180) {
		c.JSON(RespWithMsg(InvalidParam, "仅支持查询180天以内的数据"))
		return
	}
	list, err := h.service.AnalysisWeeklyTrend(c, &proto.DateRange{
		Begin: r.Begin,
		End:   r.End,
	})
	if err != nil {
		logger.FromContext(c).Error("service.AnalysisWeeklyTrend error", &r, err)
		c.JSON(RespWrapErr(err))
		return
	}
	c.JSON(OK, &AnalysisResp{
		List: list,
	})
}
