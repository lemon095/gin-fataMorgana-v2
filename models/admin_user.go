package models

import (
	"time"
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
	Role         int        `json:"role" gorm:"not null;default:1;comment:身份角色 1:超级管理员 2:经理 3:主管 4:业务员"`
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
	RoleSuperAdmin = 1 // 超级管理员
	RoleManager    = 2 // 经理
	RoleSupervisor = 3 // 主管
	RoleSalesman   = 4 // 业务员
)

// 角色名称映射
var RoleNames = map[int]string{
	RoleSuperAdmin: "超级管理员",
	RoleManager:    "经理",
	RoleSupervisor: "主管",
	RoleSalesman:   "业务员",
}

// IsActive 检查管理员是否激活
func (a *AdminUser) IsActive() bool {
	return a.Status == 1
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

// GetRoleIDByName 根据角色名称获取角色ID
func GetRoleIDByName(roleName string) (int, bool) {
	for roleID, name := range RoleNames {
		if name == roleName {
			return roleID, true
		}
	}
	return 0, false
}
