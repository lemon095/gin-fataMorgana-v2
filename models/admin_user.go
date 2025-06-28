package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AdminUser 邀请码管理表（仅用于邀请码校验）
type AdminUser struct {
	ID           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	AdminID      string     `json:"admin_id" gorm:"uniqueIndex;not null;size:8;comment:管理员唯一ID"`
	Username     string     `json:"username" gorm:"not null;size:50;uniqueIndex;comment:用户名"`
	Password     string     `json:"-" gorm:"not null;size:255;comment:密码哈希"`
	Remark       string     `json:"remark" gorm:"size:500;comment:备注"`
	Status       int        `json:"status" gorm:"default:1;comment:账户状态 1:正常 0:禁用"`
	Avatar       string     `json:"avatar" gorm:"size:255;comment:头像URL"`
	Role         int        `json:"role" gorm:"not null;default:1;comment:身份角色 1:业务员 2:主管 3:经理 4:超级管理员"`
	MyInviteCode string     `json:"my_invite_code" gorm:"size:6;uniqueIndex;comment:我的邀请码"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    *time.Time `json:"-" gorm:"index;comment:软删除时间"`
}

// TableName 指定表名
func (AdminUser) TableName() string {
	return "admin_users"
}

// TableComment 表注释
func (AdminUser) TableComment() string {
	return "邀请码管理表 - 存储邀请码信息，用于用户注册时的邀请码校验"
}

// 管理员角色常量（使用int枚举）
const (
	RoleSalesman   = 1 // 业务员
	RoleSupervisor = 2 // 主管
	RoleManager    = 3 // 经理
	RoleSuperAdmin = 4 // 超级管理员
)

// 角色名称映射
var RoleNames = map[int]string{
	RoleSalesman:   "业务员",
	RoleSupervisor: "主管",
	RoleManager:    "经理",
	RoleSuperAdmin: "超级管理员",
}

// 角色权限等级映射
var RoleLevels = map[int]int{
	RoleSalesman:   1,
	RoleSupervisor: 2,
	RoleManager:    3,
	RoleSuperAdmin: 4,
}

// AdminUserRegisterRequest 管理员注册请求
type AdminUserRegisterRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=50"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	Remark          string `json:"remark" binding:"max=500"`
	Role            int    `json:"role" binding:"required,min=1,max=4"` // 1-4对应不同角色
	InviteCode      string `json:"invite_code" binding:"required"`      // 注册时使用的邀请码
}

// AdminUserLoginRequest 管理员登录请求
type AdminUserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AdminUserResponse 管理员响应
type AdminUserResponse struct {
	ID           uint      `json:"id"`
	AdminID      string    `json:"admin_id"`
	Username     string    `json:"username"`
	Remark       string    `json:"remark"`
	Status       int       `json:"status"`
	Avatar       string    `json:"avatar"`
	Role         string    `json:"role"`         // 返回角色名称
	RoleID       int       `json:"role_id"`      // 返回角色ID
	MyInviteCode string    `json:"my_invite_code"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// AdminUserUpdateRequest 管理员更新请求
type AdminUserUpdateRequest struct {
	Username string `json:"username" binding:"min=3,max=50"`
	Remark   string `json:"remark" binding:"max=500"`
	Status   *int   `json:"status" binding:"oneof=0 1"`
	Avatar   string `json:"avatar"`
	Role     int    `json:"role" binding:"min=1,max=4"` // 1-4对应不同角色
}

// AdminUserChangePasswordRequest 管理员修改密码请求
type AdminUserChangePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// HashPassword 加密密码
func (a *AdminUser) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.Password = string(hashedPassword)
	return nil
}

// CheckPassword 验证密码
func (a *AdminUser) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	return err == nil
}

// ToResponse 转换为响应格式
func (a *AdminUser) ToResponse() AdminUserResponse {
	return AdminUserResponse{
		ID:           a.ID,
		AdminID:      a.AdminID,
		Username:     a.Username,
		Remark:       a.Remark,
		Status:       a.Status,
		Avatar:       a.Avatar,
		Role:         RoleNames[a.Role],
		RoleID:       a.Role,
		MyInviteCode: a.MyInviteCode,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}

// IsActive 检查管理员是否激活
func (a *AdminUser) IsActive() bool {
	return a.Status == 1
}

// IsSuperAdmin 检查是否为超级管理员
func (a *AdminUser) IsSuperAdmin() bool {
	return a.Role == RoleSuperAdmin
}

// IsManager 检查是否为经理
func (a *AdminUser) IsManager() bool {
	return a.Role == RoleManager
}

// IsSupervisor 检查是否为主管
func (a *AdminUser) IsSupervisor() bool {
	return a.Role == RoleSupervisor
}

// IsSalesman 检查是否为业务员
func (a *AdminUser) IsSalesman() bool {
	return a.Role == RoleSalesman
}

// HasPermission 检查是否有指定权限（使用角色ID）
func (a *AdminUser) HasPermission(requiredRoleID int) bool {
	return a.Role >= requiredRoleID
}

// HasPermissionByName 检查是否有指定权限（使用角色名称）
func (a *AdminUser) HasPermissionByName(requiredRoleName string) bool {
	// 根据角色名称获取角色ID
	for roleID, roleName := range RoleNames {
		if roleName == requiredRoleName {
			return a.HasPermission(roleID)
		}
	}
	return false
}

// GetRoleLevel 获取角色等级
func (a *AdminUser) GetRoleLevel() int {
	return RoleLevels[a.Role]
}

// GetRoleName 获取角色名称
func (a *AdminUser) GetRoleName() string {
	return RoleNames[a.Role]
}

// ValidateRoleID 验证角色ID是否有效
func ValidateRoleID(roleID int) bool {
	_, exists := RoleNames[roleID]
	return exists
}

// ValidateRole 验证角色名称是否有效
func ValidateRole(roleName string) bool {
	for _, name := range RoleNames {
		if name == roleName {
			return true
		}
	}
	return false
}

// GetRoleIDByName 根据角色名称获取角色ID
func GetRoleIDByName(roleName string) (int, bool) {
	for roleID, name := range RoleNames {
		if name == roleName {
			return roleID, true
		}
	}
	return 0, false
}
