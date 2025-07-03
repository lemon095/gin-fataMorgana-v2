package database

import (
	"context"
	"gin-fataMorgana/models"

	"gorm.io/gorm"
)

// MemberLevelRepository 用户等级配置表Repository
type MemberLevelRepository struct {
	db *gorm.DB
}

// NewMemberLevelRepository 创建用户等级配置表Repository实例
func NewMemberLevelRepository(db *gorm.DB) *MemberLevelRepository {
	return &MemberLevelRepository{db: db}
}

// GetByLevel 根据等级获取配置
func (r *MemberLevelRepository) GetByLevel(ctx context.Context, level int) (*models.MemberLevel, error) {
	var memberLevel models.MemberLevel
	err := r.db.WithContext(ctx).Where("level = ? AND status = 1", level).First(&memberLevel).Error
	if err != nil {
		return nil, err
	}
	return &memberLevel, nil
}

// GetByExperience 根据经验值获取等级配置
func (r *MemberLevelRepository) GetByExperience(ctx context.Context, experience int) (*models.MemberLevel, error) {
	var memberLevel models.MemberLevel
	err := r.db.WithContext(ctx).
		Where("min_experience <= ? AND max_experience >= ? AND status = 1", experience, experience).
		First(&memberLevel).Error
	if err != nil {
		return nil, err
	}
	return &memberLevel, nil
}

// GetAllActive 获取所有启用的等级配置
func (r *MemberLevelRepository) GetAllActive(ctx context.Context) ([]models.MemberLevel, error) {
	var memberLevels []models.MemberLevel
	err := r.db.WithContext(ctx).
		Where("status = 1").
		Order("level ASC").
		Find(&memberLevels).Error
	if err != nil {
		return nil, err
	}
	return memberLevels, nil
}

// GetLevelsByExperienceRange 根据经验值范围获取等级配置
func (r *MemberLevelRepository) GetLevelsByExperienceRange(ctx context.Context, minExp, maxExp int) ([]models.MemberLevel, error) {
	var memberLevels []models.MemberLevel
	err := r.db.WithContext(ctx).
		Where("(min_experience <= ? AND max_experience >= ?) OR (min_experience <= ? AND max_experience >= ?) AND status = 1",
			maxExp, minExp, minExp, maxExp).
		Order("level ASC").
		Find(&memberLevels).Error
	if err != nil {
		return nil, err
	}
	return memberLevels, nil
}

// GetNextLevel 获取下一个等级配置
func (r *MemberLevelRepository) GetNextLevel(ctx context.Context, currentLevel int) (*models.MemberLevel, error) {
	var memberLevel models.MemberLevel
	err := r.db.WithContext(ctx).
		Where("level > ? AND status = 1", currentLevel).
		Order("level ASC").
		First(&memberLevel).Error
	if err != nil {
		return nil, err
	}
	return &memberLevel, nil
}

// GetPreviousLevel 获取上一个等级配置
func (r *MemberLevelRepository) GetPreviousLevel(ctx context.Context, currentLevel int) (*models.MemberLevel, error) {
	var memberLevel models.MemberLevel
	err := r.db.WithContext(ctx).
		Where("level < ? AND status = 1", currentLevel).
		Order("level DESC").
		First(&memberLevel).Error
	if err != nil {
		return nil, err
	}
	return &memberLevel, nil
}

// GetMaxLevel 获取最大等级配置
func (r *MemberLevelRepository) GetMaxLevel(ctx context.Context) (*models.MemberLevel, error) {
	var memberLevel models.MemberLevel
	err := r.db.WithContext(ctx).
		Where("status = 1").
		Order("level DESC").
		First(&memberLevel).Error
	if err != nil {
		return nil, err
	}
	return &memberLevel, nil
}
