package database

import (
	"context"
	"gin-fataMorgana/models"

	"gorm.io/gorm"
)

// AmountConfigRepository 金额配置Repository
type AmountConfigRepository struct {
	db *gorm.DB
}

// NewAmountConfigRepository 创建金额配置Repository实例
func NewAmountConfigRepository() *AmountConfigRepository {
	return &AmountConfigRepository{
		db: DB,
	}
}

// GetAmountConfigsByType 根据类型获取金额配置列表
func (r *AmountConfigRepository) GetAmountConfigsByType(ctx context.Context, configType string) ([]models.AmountConfig, error) {
	var configs []models.AmountConfig

	err := r.db.WithContext(ctx).
		Where("type = ? AND is_active = ?", configType, true).
		Order("sort_order ASC, amount ASC").
		Find(&configs).Error

	return configs, err
}

// GetAmountConfigByID 根据ID获取金额配置
func (r *AmountConfigRepository) GetAmountConfigByID(ctx context.Context, id int64) (*models.AmountConfig, error) {
	var config models.AmountConfig

	err := r.db.WithContext(ctx).
		Where("id = ? AND is_active = ?", id, true).
		First(&config).Error

	if err != nil {
		return nil, err
	}

	return &config, nil
}

// CreateAmountConfig 创建金额配置
func (r *AmountConfigRepository) CreateAmountConfig(ctx context.Context, config *models.AmountConfig) error {
	return r.db.WithContext(ctx).Create(config).Error
}

// UpdateAmountConfig 更新金额配置
func (r *AmountConfigRepository) UpdateAmountConfig(ctx context.Context, config *models.AmountConfig) error {
	return r.db.WithContext(ctx).Save(config).Error
}

// DeleteAmountConfig 删除金额配置
func (r *AmountConfigRepository) DeleteAmountConfig(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.AmountConfig{}, id).Error
}

// GetAmountConfigsByTypeAndAmount 根据类型和金额获取配置
func (r *AmountConfigRepository) GetAmountConfigsByTypeAndAmount(ctx context.Context, configType string, amount float64) (*models.AmountConfig, error) {
	var config models.AmountConfig

	err := r.db.WithContext(ctx).
		Where("type = ? AND amount = ? AND is_active = ?", configType, amount, true).
		First(&config).Error

	if err != nil {
		return nil, err
	}

	return &config, nil
}
