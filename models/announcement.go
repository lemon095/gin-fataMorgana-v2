package models

import (
	"time"

	"gorm.io/gorm"
)

// Announcement 公告表
type Announcement struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string         `json:"title" gorm:"not null;size:100;comment:标题"`
	Content     string         `json:"content" gorm:"type:text;not null;comment:纯文本内容（用于摘要、搜索等）"`
	RichContent *string        `json:"rich_content" gorm:"type:longtext;comment:富文本内容（HTML格式）"`
	Tag         string         `json:"tag" gorm:"not null;size:20;comment:标签"`
	Status      int64          `json:"status" gorm:"default:0;comment:状态 0-草稿 1-已发布"`
	IsPublish   bool           `json:"is_publish" gorm:"default:false;comment:是否发布"`
	CreatedAt   *time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   *time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;comment:软删除时间"`

	// 关联字段
	Banners []AnnouncementBanner `json:"banners" gorm:"foreignKey:AnnouncementID"`
}

// TableName 指定表名
func (Announcement) TableName() string {
	return "announcements"
}

// TableComment 表注释
func (Announcement) TableComment() string {
	return "公告表 - 存储系统公告信息，支持富文本内容，包括标题、纯文本内容、富文本内容、标签、状态等"
}

// AnnouncementBanner 公告图片表
type AnnouncementBanner struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	AnnouncementID uint           `json:"announcement_id" gorm:"not null;comment:公告ID"`
	ImageURL       string         `json:"image_url" gorm:"not null;size:255;comment:图片URL"`
	Title          string         `json:"title" gorm:"size:100;comment:图片标题"`
	Link           string         `json:"link" gorm:"size:255;comment:跳转链接"`
	Sort           int64          `json:"sort" gorm:"default:0;comment:排序"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index;comment:软删除时间"`
}

// TableName 指定表名
func (AnnouncementBanner) TableName() string {
	return "announcement_banners"
}

// TableComment 表注释
func (AnnouncementBanner) TableComment() string {
	return "公告图片表 - 存储公告相关的图片信息"
}

// AnnouncementListRequest 公告列表请求
type AnnouncementListRequest struct {
	Page     int `json:"page" binding:"min=1"`              // 页码，从1开始
	PageSize int `json:"page_size" binding:"min=1,max=20"` // 每页大小，最大20
}

// AnnouncementResponse 公告响应
type AnnouncementResponse struct {
	ID          uint             `json:"id"`
	Title       string           `json:"title"`
	Content     string           `json:"content"`
	RichContent *string          `json:"rich_content"`
	Tag         string           `json:"tag"`
	Status      int64            `json:"status"`
	IsPublish   bool             `json:"is_publish"`
	CreatedAt   *time.Time       `json:"created_at"`
	Banners     []BannerResponse `json:"banners"`
}

// BannerResponse 图片响应（直接返回图片URL字符串）
type BannerResponse string

// AnnouncementListResponse 公告列表响应
type AnnouncementListResponse struct {
	Announcements []AnnouncementResponse `json:"announcements"`
	Pagination    PaginationInfo         `json:"pagination"`
}

// ToResponse 转换为响应格式
func (a *Announcement) ToResponse() AnnouncementResponse {
	var banners []BannerResponse
	for _, banner := range a.Banners {
		banners = append(banners, BannerResponse(banner.ImageURL))
	}

	return AnnouncementResponse{
		ID:          a.ID,
		Title:       a.Title,
		Content:     a.Content,
		RichContent: a.RichContent,
		Tag:         a.Tag,
		Status:      a.Status,
		IsPublish:   a.IsPublish,
		CreatedAt:   a.CreatedAt,
		Banners:     banners,
	}
}

// IsPublished 检查公告是否已发布
func (a *Announcement) IsPublished() bool {
	return a.Status == 1 && a.IsPublish
}
