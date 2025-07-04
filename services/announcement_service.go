package services

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
	"math"
	"time"
)

// AnnouncementService 公告服务
type AnnouncementService struct {
	repo *database.AnnouncementRepository
}

// NewAnnouncementService 创建公告服务实例
func NewAnnouncementService() *AnnouncementService {
	return &AnnouncementService{
		repo: database.NewAnnouncementRepository(),
	}
}

// GetAnnouncementList 获取公告列表（带缓存）
func (s *AnnouncementService) GetAnnouncementList(ctx context.Context, page, pageSize int) (*models.AnnouncementListResponse, error) {
	// 验证分页参数
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 20 {
		return nil, fmt.Errorf("每页数量不能超过20")
	}

	// 生成缓存键
	cacheKey := fmt.Sprintf("announcement:list:page:%d:size:%d", page, pageSize)

	// 尝试从缓存获取数据
	if cached, err := database.RedisClient.Get(ctx, cacheKey).Result(); err == nil {
		var response models.AnnouncementListResponse
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			return &response, nil
		}
	}

	// 缓存未命中，从数据库获取数据
	announcements, total, err := s.repo.GetAnnouncementList(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	var responses []models.AnnouncementResponse
	for _, announcement := range announcements {
		responses = append(responses, announcement.ToResponse())
	}

	// 计算分页信息
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	hasNext := page < totalPages
	hasPrev := page > 1

	pagination := models.PaginationInfo{
		CurrentPage: page,
		PageSize:    pageSize,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
	}

	response := &models.AnnouncementListResponse{
		Announcements: responses,
		Pagination:    pagination,
	}

	// 缓存数据（1分钟）
	if data, err := json.Marshal(response); err == nil {
		database.RedisClient.Set(ctx, cacheKey, data, time.Minute)
	}

	return response, nil
}

// GetAnnouncementByID 根据ID获取公告详情
func (s *AnnouncementService) GetAnnouncementByID(ctx context.Context, id uint) (*models.AnnouncementResponse, error) {
	announcement, err := s.repo.GetAnnouncementByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := announcement.ToResponse()
	return &response, nil
}

// ClearAnnouncementCache 清除公告缓存
func (s *AnnouncementService) ClearAnnouncementCache(ctx context.Context) error {
	pattern := "announcement:list:*"
	keys, err := database.RedisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return database.RedisClient.Del(ctx, keys...).Err()
	}
	return nil
}
