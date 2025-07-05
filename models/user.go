package models

import (
	"gin-fataMorgana/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// 用户状态常量
const (
	UserStatusDisabled = 0 // 禁用
	UserStatusActive   = 1 // 正常
	UserStatusPending  = 2 // 待审核
)

// User 用户模型
type User struct {
	ID                       uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Uid                      string     `json:"uid" gorm:"uniqueIndex;not null;size:8;comment:用户唯一ID"`
	Username                 string     `json:"username" gorm:"not null;size:50;index;comment:用户名"`
	Email                    string     `json:"email" gorm:"uniqueIndex;not null;size:100;comment:邮箱地址"`
	Password                 string     `json:"-" gorm:"not null;size:255;comment:密码哈希"`
	Phone                    string     `json:"phone" gorm:"size:20;index;comment:手机号"`
	BankCardInfo             string     `json:"bank_card_info" gorm:"type:json;comment:银行卡信息JSON"`
	Experience               int        `json:"experience" gorm:"default:0;comment:用户经验值"`
	CreditScore              int        `json:"credit_score" gorm:"default:100;comment:用户信用分"`
	Status                   int        `json:"status" gorm:"default:2;comment:用户状态 0:禁用 1:正常 2:待审核"`
	InvitedBy                string     `json:"invited_by" gorm:"size:6;index;comment:注册时填写的邀请码"`
	HasGroupBuyQualification bool       `json:"has_group_buy_qualification" gorm:"default:false;comment:是否有拼单资格"`
	Rate                     int        `json:"rate" gorm:"-"` // 用户等级进度（从Redis获取，不存储到数据库）
	CreatedAt                time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt                time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt                *time.Time `json:"-" gorm:"index;comment:软删除时间"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// TableComment 表注释
func (User) TableComment() string {
	return "用户表 - 存储用户基本信息、认证信息、银行卡信息、经验值、信用分等"
}

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Account         string `json:"account" binding:"required"` // 可以是邮箱或手机号
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	InviteCode      string `json:"invite_code" binding:"required"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Account  string `json:"account" binding:"required"` // 可以是邮箱或手机号
	Password string `json:"password" binding:"required"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID                       uint      `json:"id"`
	Uid                      string    `json:"uid"`
	Username                 string    `json:"username"`
	Email                    string    `json:"email"`
	Phone                    string    `json:"phone"`
	BankCardInfo             string    `json:"bank_card_info"`
	Experience               int       `json:"experience"`
	CreditScore              int       `json:"credit_score"`
	Status                   int       `json:"status"`
	InvitedBy                string    `json:"invited_by"`
	HasGroupBuyQualification bool      `json:"has_group_buy_qualification"`
	Rate                     int       `json:"rate"` // 用户等级进度（整数）
	CreatedAt                time.Time `json:"created_at"`
}

// TokenResponse Token响应
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// RefreshTokenRequest 刷新Token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// BankCardInfo 银行卡信息结构
type BankCardInfo struct {
	CardNumber string `json:"card_number"`
	CardType   string `json:"card_type"` // 借记卡、信用卡等
	BankName   string `json:"bank_name"`
	CardHolder string `json:"card_holder"`
}

// GetProfileRequest 获取用户信息请求
type GetProfileRequest struct {
	// 空结构体，因为获取当前用户信息不需要额外参数
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`       // 旧密码
	NewPassword string `json:"new_password" binding:"required,min=6"` // 新密码
}

// GetBankCardRequest 获取银行卡信息请求
type GetBankCardRequest struct {
	// 空结构体，因为获取当前用户银行卡信息不需要额外参数
}

// GetSessionStatusRequest 获取会话状态请求
type GetSessionStatusRequest struct {
	// 空结构体，因为获取会话状态不需要额外参数
}

// GetSessionUserRequest 获取会话用户信息请求
type GetSessionUserRequest struct {
	// 空结构体，因为获取会话用户信息不需要额外参数
}

// GetWalletRequest 获取钱包信息请求
type GetWalletRequest struct {
	// 空结构体，因为获取当前用户钱包信息不需要额外参数
}

// GetTransactionsRequest 获取交易记录请求
type GetTransactionsRequest struct {
	Page     int    `json:"page" binding:"min=1"`              // 页码，从1开始
	PageSize int    `json:"page_size" binding:"min=1"`         // 每页大小，最小1
	Type     string `json:"type"`                              // 交易类型过滤：recharge(充值)、withdraw(提现)、order_buy(购买)、group_buy(拼单)
}

// GetWithdrawSummaryRequest 获取提现汇总请求
type GetWithdrawSummaryRequest struct {
	// 空结构体，因为获取提现汇总不需要额外参数
}

// HealthCheckRequest 健康检查请求
type HealthCheckRequest struct {
	// 空结构体，健康检查不需要参数
}

// HashPassword 加密密码
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// ToResponse 转换为响应格式
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:                       u.ID,
		Uid:                      u.Uid,
		Username:                 u.Username,
		Email:                    utils.MaskEmail(u.Email),
		Phone:                    utils.MaskPhone(u.Phone),
		BankCardInfo:             u.BankCardInfo,
		Experience:               u.Experience,
		CreditScore:              u.CreditScore,
		Status:                   u.Status,
		InvitedBy:                u.InvitedBy,
		HasGroupBuyQualification: u.HasGroupBuyQualification,
		Rate:                     0, // 默认值，实际值需要从Redis获取
		CreatedAt:                u.CreatedAt,
	}
}

// IsActive 检查用户是否激活
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// CheckGroupBuyQualification 检查用户是否有拼单资格
func (u *User) CheckGroupBuyQualification() bool {
	return u.HasGroupBuyQualification
}

// SetGroupBuyQualification 设置用户拼单资格
func (u *User) SetGroupBuyQualification(hasQualification bool) {
	u.HasGroupBuyQualification = hasQualification
}
