package wechat

import "context"

/*
数据分析：
GetDailySummary：日访问概况
GetDailyVisitTrend：日访问趋势
GetWeeklyVisitTrend：周访问趋势
*/

const DateFormat = "20060102"

type DatacubeArgs struct {
	BeginDate string `json:"begin_date"` // yyyymmdd
	EndDate   string `json:"end_date"`   // yyyymmdd
}

type GetDailySummaryResp struct {
	respErr
	List []*SummaryData `json:"list,omitempty"`
}
type SummaryData struct {
	RefDate    string `json:"ref_date"` // yyyymmdd
	VisitTotal int    `json:"visit_total"`
	SharePv    int    `json:"share_pv"`
	ShareUv    int    `json:"share_uv"`
}

func (api *server) GetDailySummary(ctx context.Context, args *DatacubeArgs) (*GetDailySummaryResp, error) {
	var resp GetDailySummaryResp
	err := api.post(ctx, "/datacube/getweanalysisappiddailysummarytrend", args, &resp)
	return &resp, err
}

type VisitTrendResp struct {
	respErr
	List []*VisitTrendData `json:"list"`
}
type VisitTrendData struct {
	RefDate         string  `json:"ref_date"`
	SessionCnt      int     `json:"session_cnt"`
	VisitPv         int     `json:"visit_pv"`
	VisitUv         int     `json:"visit_uv"`
	VisitUvNew      int     `json:"visit_uv_new"`
	StayTimeUv      float64 `json:"stay_time_uv"`
	StayTimeSession float64 `json:"stay_time_session"`
	VisitDepth      float64 `json:"visit_depth"`
}

func (api *server) GetDailyVisitTrend(ctx context.Context, args *DatacubeArgs) (*VisitTrendResp, error) {
	var resp VisitTrendResp
	err := api.post(ctx, "/datacube/getweanalysisappiddailyvisittrend", args, &resp)
	return &resp, err
}

func (api *server) GetWeeklyVisitTrend(ctx context.Context, args *DatacubeArgs) (*VisitTrendResp, error) {
	var resp VisitTrendResp
	err := api.post(ctx, "/datacube/getweanalysisappidweeklyvisittrend", args, &resp)
	return &resp, err
}
