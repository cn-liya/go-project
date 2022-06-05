package model

type AnalysisDailySummary struct {
	RefDate    string `json:"ref_date"`
	VisitTotal int    `json:"visit_total"` //累计用户数
	SharePv    int    `json:"share_pv"`    //转发次数
	ShareUv    int    `json:"share_uv"`    //转发人数
}

func (*AnalysisDailySummary) TableName() string {
	return "analysis_daily_summary"
}

type visitTrend struct {
	RefDate         string  `json:"ref_date"`
	SessionCnt      int     `json:"session_cnt"`       //打开次数
	VisitPv         int     `json:"visit_pv"`          //访问次数
	VisitUv         int     `json:"visit_uv"`          //访问人数
	VisitUvNew      int     `json:"visit_uv_new"`      //新用户数
	StayTimeUv      float64 `json:"stay_time_uv"`      //人均停留时长
	StayTimeSession float64 `json:"stay_time_session"` //次均停留时长
	VisitDepth      float64 `json:"visit_depth"`       //平均访问深度
}

type AnalysisDailyTrend visitTrend

func (*AnalysisDailyTrend) TableName() string {
	return "analysis_daily_trend"
}

type AnalysisWeeklyTrend visitTrend

func (*AnalysisWeeklyTrend) TableName() string {
	return "analysis_weekly_trend"
}
