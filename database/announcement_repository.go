package database

import (
	"context"
	"gin-fataMorgana/models"

	"gorm.io/gorm"
)

// AnnouncementRepository 公告Repository
type AnnouncementRepository struct {
	db *gorm.DB
}

// NewAnnouncementRepository 创建公告Repository实例
func NewAnnouncementRepository() *AnnouncementRepository {
	return &AnnouncementRepository{
		db: DB,
	}
}

// GetAnnouncementList 获取公告列表
func (r *AnnouncementRepository) GetAnnouncementList(ctx context.Context, page, pageSize int) ([]models.Announcement, int64, error) {
	var announcements []models.Announcement
	var total int64
	
	// 计算偏移量
	offset := (page - 1) * pageSize
	
	// 获取总数
	err := r.db.WithContext(ctx).
		Model(&models.Announcement{}).
		Where("status = ? AND is_publish = ?", 1, true).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// 获取公告列表，按创建时间倒序排列
	err = r.db.WithContext(ctx).
		Preload("Banners", func(db *gorm.DB) *gorm.DB {
			return db.Select("announcement_id, image_url").Order("sort ASC")
		}).
		Where("status = ? AND is_publish = ?", 1, true).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&announcements).Error
	
	return announcements, total, err
}

// GetAnnouncementByID 根据ID获取公告详情
func (r *AnnouncementRepository) GetAnnouncementByID(ctx context.Context, id uint) (*models.Announcement, error) {
	var announcement models.Announcement
	
	err := r.db.WithContext(ctx).
		Preload("Banners", func(db *gorm.DB) *gorm.DB {
			return db.Select("announcement_id, image_url").Order("sort ASC")
		}).
		Where("id = ? AND status = ? AND is_publish = ?", id, 1, true).
		First(&announcement).Error
	
	if err != nil {
		return nil, err
	}
	
	return &announcement, nil
} 