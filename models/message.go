package models

import (
	"time"
)

// Message 消息表
type Message struct {
	ID          int64      `json:"id" gorm:"primaryKey;autoIncrement;comment:唯一自增ID"`
	UID         string     `json:"uid" gorm:"not null;index;size:8;comment:用户唯一ID"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime;type:datetime(3);comment:创建时间"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime;type:datetime(3);comment:更新时间"`
	Status      string     `json:"status" gorm:"not null;default:draft;type:enum('draft','sent','read');comment:状态：draft-草稿/未发送，sent-已发送，read-已读"`
	MessageType string     `json:"message_type" gorm:"not null;default:info;type:enum('warning','error','question','info');comment:消息类型：warning-警告，error-错误，question-问号，info-消息"`
	Content     string     `json:"content" gorm:"not null;type:text;comment:消息内容"`
	ReadAt      *time.Time `json:"read_at" gorm:"type:datetime(3);comment:用户读取时间"`
	CreatedBy   string     `json:"created_by" gorm:"not null;size:100;comment:创建管理员用户名"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "messages"
}

// UserMessage Redis消息结构
type UserMessage struct {
	ID          int64  `json:"id"`
	MessageType string `json:"message_type"`
	Content     string `json:"content"`
	CreatedAt   string `json:"created_at"`
}

// UserMessageResponse 用户消息响应
type UserMessageResponse struct {
	MessageType string `json:"message_type"`
	Content     string `json:"content"`
}