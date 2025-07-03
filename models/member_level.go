package models

import (
	"time"

	"gorm.io/gorm"
)

// MemberLevel 用户等级配置表
type MemberLevel struct {
	ID            uint64         `gorm:"primarykey" json:"id"`
	Level         int            `gorm:"not null;uniqueIndex:uniq_level;comment:等级数值" json:"level"`
	Name          string         `gorm:"size:20;not null;comment:等级名称" json:"name"`
	Logo          string         `gorm:"size:255;comment:等级logo" json:"logo"`
	Remark        string         `gorm:"size:255;comment:备注" json:"remark"`
	CashbackRatio float64        `gorm:"type:decimal(5,2);default:0;comment:返现比例（百分比）" json:"cashback_ratio"`
	SingleAmount  int            `gorm:"default:1;comment:单数字额" json:"single_amount"`
	CreatedAt     time.Time      `gorm:"type:datetime(3);autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"type:datetime(3);autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"type:datetime(3);index;comment:软删除时间" json:"-"`
}

// TableName 指定表名
func (MemberLevel) TableName() string {
	return "member_level"
}

// TableComment 表注释
func (MemberLevel) TableComment() string {
	return "用户等级配置表 - 存储用户等级配置信息，包括等级、名称、logo、返现比例、单数字额等"
}

// GetCashbackRatio 获取返现比例
func (ml *MemberLevel) GetCashbackRatio() float64 {
	return ml.CashbackRatio
}

// GetSingleAmount 获取单数字额
func (ml *MemberLevel) GetSingleAmount() int {
	return ml.SingleAmount
}

// IsDeleted 检查是否已软删除
func (ml *MemberLevel) IsDeleted() bool {
	return ml.DeletedAt.Valid
}

// MemberLevelRequest 等级配置请求
type MemberLevelRequest struct {
	Level         int     `json:"level" binding:"required,min=1"`
	Name          string  `json:"name" binding:"required,max=20"`
	Logo          string  `json:"logo"`
	Remark        string  `json:"remark"`
	CashbackRatio float64 `json:"cashback_ratio" binding:"min=0,max=100"`
	SingleAmount  int     `json:"single_amount" binding:"min=1"`
}

// MemberLevelResponse 等级配置响应
type MemberLevelResponse struct {
	ID            uint64    `json:"id"`
	Level         int       `json:"level"`
	Name          string    `json:"name"`
	Logo          string    `json:"logo"`
	Remark        string    `json:"remark"`
	CashbackRatio float64   `json:"cashback_ratio"`
	SingleAmount  int       `json:"single_amount"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ToResponse 转换为响应格式
func (ml *MemberLevel) ToResponse() MemberLevelResponse {
	return MemberLevelResponse{
		ID:            ml.ID,
		Level:         ml.Level,
		Name:          ml.Name,
		Logo:          ml.Logo,
		Remark:        ml.Remark,
		CashbackRatio: ml.CashbackRatio,
		SingleAmount:  ml.SingleAmount,
		CreatedAt:     ml.CreatedAt,
		UpdatedAt:     ml.UpdatedAt,
	}
}
