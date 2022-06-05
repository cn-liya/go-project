package service

import (
	"context"
	"project/cms/internal/proto"
	"project/model"
)

func (s *Service) AnalysisDailySummary(ctx context.Context, p *proto.DateRange) ([]*model.AnalysisDailySummary, error) {
	var list []*model.AnalysisDailySummary
	err := s.mysql.WithContext(ctx).
		Where("ref_date >= ? AND ref_date <= ?", p.Begin, p.End).
		Order("ref_date").Find(&list).Error
	return list, err
}

func (s *Service) AnalysisDailyTrend(ctx context.Context, p *proto.DateRange) ([]*model.AnalysisDailyTrend, error) {
	var list []*model.AnalysisDailyTrend
	err := s.mysql.WithContext(ctx).
		Where("ref_date >= ? AND ref_date <= ?", p.Begin, p.End).
		Order("ref_date").Find(&list).Error
	return list, err
}

func (s *Service) AnalysisWeeklyTrend(ctx context.Context, p *proto.DateRange) ([]*model.AnalysisWeeklyTrend, error) {
	var list []*model.AnalysisWeeklyTrend
	err := s.mysql.WithContext(ctx).
		Where("ref_date >= ? AND ref_date <= ?", p.Begin, p.End).
		Order("ref_date").Find(&list).Error
	return list, err
}
