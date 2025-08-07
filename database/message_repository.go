package database

import (
	"context"
	"gin-fataMorgana/models"
	"time"

	"gorm.io/gorm"
)

// MessageRepository 消息Repository
type MessageRepository struct {
	db *gorm.DB
}

// NewMessageRepository 创建消息Repository实例
func NewMessageRepository() *MessageRepository {
	return &MessageRepository{
		db: DB,
	}
}

// GetMessageByID 根据ID获取消息
func (r *MessageRepository) GetMessageByID(ctx context.Context, id int64) (*models.Message, error) {
	var message models.Message

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&message).Error

	if err != nil {
		return nil, err
	}

	return &message, nil
}

// UpdateMessageStatus 更新消息状态为已读
func (r *MessageRepository) UpdateMessageStatus(ctx context.Context, id int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.Message{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     "read",
			"read_at":    &now,
			"updated_at": now,
		}).Error
}

// CreateMessage 创建消息
func (r *MessageRepository) CreateMessage(ctx context.Context, message *models.Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}

// GetUserMessages 获取用户消息列表
func (r *MessageRepository) GetUserMessages(ctx context.Context, uid string, limit int) ([]models.Message, error) {
	var messages []models.Message

	err := r.db.WithContext(ctx).
		Where("uid = ?", uid).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error

	return messages, err
}