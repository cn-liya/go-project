package service

import (
	"context"
	"encoding/json"
	"gorm.io/gorm"
	"project/cms/internal/proto"
	"project/model/entity"
)

func (s *Service) SystemActionKeyNames(ctx context.Context, level int8) ([]string, error) {
	var result []string
	err := s.mysql.WithContext(ctx).
		Raw("SELECT key_name FROM system_action WHERE level=?", level).Scan(&result).Error
	return result, err
}

func (s *Service) SystemActionList(ctx context.Context) ([]*entity.SystemAction, error) {
	var result []*entity.SystemAction
	err := s.mysql.WithContext(ctx).Find(&result).Error
	return result, err
}

func (s *Service) SystemActionUpdate(ctx context.Context, data *entity.SystemAction) error {
	return s.mysql.WithContext(ctx).Updates(data).Error
}

func (s *Service) fetchSystemConfig(ctx context.Context, keyName string, val any) error {
	var content string
	err := s.mysql.WithContext(ctx).
		Raw("SELECT content FROM system_config WHERE key_name=?", keyName).
		Scan(&content).Error
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(content), val)
	return err
}

func (s *Service) SystemMenuKeyNames(ctx context.Context) ([]string, error) {
	var result []string
	err := s.fetchSystemConfig(ctx, entity.SysCfgMenuKeys, &result)
	return result, err
}

func (s *Service) SystemMenuTrees(ctx context.Context) ([]*entity.SysMenuTree, error) {
	var result []*entity.SysMenuTree
	err := s.fetchSystemConfig(ctx, entity.SysCfgMenuTree, &result)
	return result, err
}

func (s *Service) SystemActionMenuSave(ctx context.Context, p *proto.SystemSyncData) error {
	return s.mysql.WithContext(ctx).Transaction(func(tx *gorm.DB) (err error) {
		if len(p.ActionDelete) > 0 {
			if err = tx.Delete(&entity.SystemAction{}, p.ActionDelete).Error; err != nil {
				return
			}
		}
		if len(p.ActionCreate) > 0 {
			if err = tx.Create(p.ActionCreate).Error; err != nil {
				return
			}
		}
		if err = tx.Updates(&entity.SystemConfig{
			KeyName: entity.SysCfgMenuTree,
			Content: entity.JsonString(p.MenuTree),
		}).Error; err != nil {
			return
		}
		err = tx.Updates(&entity.SystemConfig{
			KeyName: entity.SysCfgMenuKeys,
			Content: entity.JsonString(p.MenuKeys),
		}).Error
		return
	})
}
