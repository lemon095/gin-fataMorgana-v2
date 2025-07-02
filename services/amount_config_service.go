package services

import (
	"context"
	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
)

// AmountConfigService 金额配置服务
type AmountConfigService struct {
	repo *database.AmountConfigRepository
}

// NewAmountConfigService 创建金额配置服务实例
func NewAmountConfigService() *AmountConfigService {
	return &AmountConfigService{
		repo: database.NewAmountConfigRepository(),
	}
}

// GetAmountConfigsByType 根据类型获取金额配置列表
func (s *AmountConfigService) GetAmountConfigsByType(ctx context.Context, configType string) ([]*models.AmountConfigResponse, error) {
	configs, err := s.repo.GetAmountConfigsByType(ctx, configType)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	var responses []*models.AmountConfigResponse
	for _, config := range configs {
		responses = append(responses, config.ToResponse())
	}

	return responses, nil
}

// GetAmountConfigByID 根据ID获取金额配置
func (s *AmountConfigService) GetAmountConfigByID(ctx context.Context, id int64) (*models.AmountConfigResponse, error) {
	config, err := s.repo.GetAmountConfigByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return config.ToResponse(), nil
}

// CreateAmountConfig 创建金额配置
func (s *AmountConfigService) CreateAmountConfig(ctx context.Context, config *models.AmountConfig) (*models.AmountConfigResponse, error) {
	err := s.repo.CreateAmountConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return config.ToResponse(), nil
}

// UpdateAmountConfig 更新金额配置
func (s *AmountConfigService) UpdateAmountConfig(ctx context.Context, config *models.AmountConfig) (*models.AmountConfigResponse, error) {
	err := s.repo.UpdateAmountConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return config.ToResponse(), nil
}

// DeleteAmountConfig 删除金额配置
func (s *AmountConfigService) DeleteAmountConfig(ctx context.Context, id int64) error {
	return s.repo.DeleteAmountConfig(ctx, id)
}

// GetAmountConfigsByTypeAndAmount 根据类型和金额获取配置
func (s *AmountConfigService) GetAmountConfigsByTypeAndAmount(ctx context.Context, configType string, amount float64) (*models.AmountConfigResponse, error) {
	config, err := s.repo.GetAmountConfigsByTypeAndAmount(ctx, configType, amount)
	if err != nil {
		return nil, err
	}

	return config.ToResponse(), nil
} 