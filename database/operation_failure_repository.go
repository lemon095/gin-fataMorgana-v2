package database

import (
	"context"
	"gin-fataMorgana/models"
)

// OperationFailureRepository 操作失败记录Repository
type OperationFailureRepository struct {
	*BaseRepository
}

func NewOperationFailureRepository() *OperationFailureRepository {
	return &OperationFailureRepository{
		BaseRepository: NewBaseRepository(),
	}
}

// Create 创建失败记录
func (r *OperationFailureRepository) Create(ctx context.Context, failure *models.OperationFailure) error {
	return r.Create(ctx, failure)
}

// GetByID 根据ID获取失败记录
func (r *OperationFailureRepository) GetByID(ctx context.Context, id uint) (*models.OperationFailure, error) {
	var failure models.OperationFailure
	err := r.FindByCondition(ctx, map[string]interface{}{"id": id}, &failure)
	if err != nil {
		return nil, err
	}
	return &failure, nil
}

// GetByUID 获取用户的失败记录
func (r *OperationFailureRepository) GetByUID(ctx context.Context, uid string, page, pageSize int) ([]*models.OperationFailure, int64, error) {
	var failures []*models.OperationFailure
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	err := r.db.WithContext(ctx).Model(&models.OperationFailure{}).Where("uid = ?", uid).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = r.db.WithContext(ctx).Where("uid = ?", uid).
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&failures).Error

	return failures, total, err
}

// GetByOperationType 根据操作类型获取失败记录
func (r *OperationFailureRepository) GetByOperationType(ctx context.Context, operationType string, page, pageSize int) ([]*models.OperationFailure, int64, error) {
	var failures []*models.OperationFailure
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	err := r.db.WithContext(ctx).Model(&models.OperationFailure{}).Where("operation_type = ?", operationType).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = r.db.WithContext(ctx).Where("operation_type = ?", operationType).
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&failures).Error

	return failures, total, err
}

// GetRecentFailures 获取最近的失败记录
func (r *OperationFailureRepository) GetRecentFailures(ctx context.Context, limit int) ([]*models.OperationFailure, error) {
	var failures []*models.OperationFailure

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Find(&failures).Error

	return failures, err
}
