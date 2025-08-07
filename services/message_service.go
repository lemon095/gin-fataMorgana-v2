package services

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
	"gin-fataMorgana/utils"

	"github.com/redis/go-redis/v9"
)

// MessageService 消息服务
type MessageService struct {
	repo        *database.MessageRepository
	redisClient *redis.Client
}

// NewMessageService 创建消息服务实例
func NewMessageService() *MessageService {
	return &MessageService{
		repo:        database.NewMessageRepository(),
		redisClient: database.RedisClient,
	}
}

// GetUserMessage 获取用户消息推送
func (s *MessageService) GetUserMessage(ctx context.Context, uid string) (*models.UserMessageResponse, error) {
	// Redis key
	key := fmt.Sprintf("user_messages:%s", uid)

	// 获取列表长度
	length, err := s.redisClient.LLen(ctx, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			// 列表不存在，返回空
			return nil, nil
		}
		return nil, utils.NewAppError(utils.CodeRedisError, "获取消息列表失败")
	}

	// 如果列表为空，返回空
	if length == 0 {
		return nil, nil
	}

	// 从列表左侧pop一条消息
	messageJSON, err := s.redisClient.LPop(ctx, key).Result()
	if err != nil {
		return nil, utils.NewAppError(utils.CodeRedisError, "获取消息失败")
	}

	// 解析消息JSON
	var userMessage models.UserMessage
	if err := json.Unmarshal([]byte(messageJSON), &userMessage); err != nil {
		return nil, utils.NewAppError(utils.CodeServer, "解析消息数据失败")
	}

	// 更新消息状态为已读
	if err := s.repo.UpdateMessageStatus(ctx, userMessage.ID); err != nil {
		// 如果更新失败，将消息重新放回列表头部
		s.redisClient.LPush(ctx, key, messageJSON)
		return nil, utils.NewAppError(utils.CodeDatabaseError, "更新消息状态失败")
	}

	// 返回消息内容
	return &models.UserMessageResponse{
		MessageType: userMessage.MessageType,
		Content:     userMessage.Content,
	}, nil
}

// PushUserMessage 推送消息到用户
func (s *MessageService) PushUserMessage(ctx context.Context, uid string, messageType, content, createdBy string) error {
	// 创建消息记录
	message := &models.Message{
		UID:         uid,
		Status:      "sent",
		MessageType: messageType,
		Content:     content,
		CreatedBy:   createdBy,
	}

	// 保存到数据库
	if err := s.repo.CreateMessage(ctx, message); err != nil {
		return utils.NewAppError(utils.CodeDatabaseError, "创建消息记录失败")
	}

	// 构建Redis消息
	userMessage := models.UserMessage{
		ID:          message.ID,
		MessageType: message.MessageType,
		Content:     message.Content,
		CreatedAt:   message.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	// 序列化为JSON
	messageJSON, err := json.Marshal(userMessage)
	if err != nil {
		return utils.NewAppError(utils.CodeServer, "序列化消息数据失败")
	}

	// 推送到Redis列表
	key := fmt.Sprintf("user_messages:%s", uid)
	if err := s.redisClient.RPush(ctx, key, messageJSON).Err(); err != nil {
		return utils.NewAppError(utils.CodeRedisError, "推送消息失败")
	}

	return nil
}