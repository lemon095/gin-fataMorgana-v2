package services

import (
	"context"
	"gin-fataMorgana/database"
)

const ShareLinkKey = "admin_system_share_link" // 分享链接缓存键

// ShareService 分享相关服务
type ShareService struct{}

func NewShareService() *ShareService {
	return &ShareService{}
}

// GetShareLink 获取分享链接
func (s *ShareService) GetShareLink(ctx context.Context) (string, error) {
	link, err := database.GetKey(ctx, ShareLinkKey)
	if err != nil {
		return "", err
	}
	return link, nil
}
