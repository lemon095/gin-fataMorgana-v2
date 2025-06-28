package models

import (
	"time"
)

// UserLoginLog 用户登录记录表
type UserLoginLog struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Uid        string    `json:"uid" gorm:"not null;size:8;index;comment:用户UID"`
	Username   string    `json:"username" gorm:"not null;size:50;index;comment:用户名"`
	Email      string    `json:"email" gorm:"not null;size:100;index;comment:邮箱"`
	LoginIP    string    `json:"login_ip" gorm:"not null;size:45;index;comment:登录IP地址"`
	UserAgent  string    `json:"user_agent" gorm:"size:500;comment:用户代理"`
	LoginTime  time.Time `json:"login_time" gorm:"not null;index;comment:登录时间"`
	Status     int       `json:"status" gorm:"default:1;comment:登录状态 1:成功 0:失败"`
	FailReason string    `json:"fail_reason" gorm:"size:200;comment:失败原因"`
	DeviceInfo string    `json:"device_info" gorm:"size:200;comment:设备信息"`
	Location   string    `json:"location" gorm:"size:100;comment:登录地点"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName 指定表名
func (UserLoginLog) TableName() string {
	return "user_login_logs"
}

// TableComment 表注释
func (UserLoginLog) TableComment() string {
	return "用户登录日志表 - 记录用户登录历史，包括登录时间、IP地址、设备信息、登录状态等"
}

// UserLoginLogResponse 登录记录响应
type UserLoginLogResponse struct {
	ID         uint      `json:"id"`
	Uid        string    `json:"uid"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	LoginIP    string    `json:"login_ip"`
	UserAgent  string    `json:"user_agent"`
	LoginTime  time.Time `json:"login_time"`
	Status     int       `json:"status"`
	FailReason string    `json:"fail_reason"`
	DeviceInfo string    `json:"device_info"`
	Location   string    `json:"location"`
	CreatedAt  time.Time `json:"created_at"`
}

// ToResponse 转换为响应格式
func (l *UserLoginLog) ToResponse() UserLoginLogResponse {
	return UserLoginLogResponse{
		ID:         l.ID,
		Uid:        l.Uid,
		Username:   l.Username,
		Email:      l.Email,
		LoginIP:    l.LoginIP,
		UserAgent:  l.UserAgent,
		LoginTime:  l.LoginTime,
		Status:     l.Status,
		FailReason: l.FailReason,
		DeviceInfo: l.DeviceInfo,
		Location:   l.Location,
		CreatedAt:  l.CreatedAt,
	}
}

// IsSuccess 检查登录是否成功
func (l *UserLoginLog) IsSuccess() bool {
	return l.Status == 1
}
