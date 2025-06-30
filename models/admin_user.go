package models

import (
	"time"

	"gorm.io/gorm"
)

// AdminUser 邀请码管理表（仅用于邀请码校验）
type AdminUser struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	AdminID      uint           `gorm:"not null;uniqueIndex:idx_admin_users_admin_id;comment:管理员唯一ID" json:"admin_id"`
	Username     string         `gorm:"size:50;not null;uniqueIndex:idx_admin_users_username;comment:用户名" json:"username"`
	Password     string         `gorm:"size:255;not null;comment:密码哈希" json:"-"`
	Remark       string         `gorm:"size:500;comment:备注" json:"remark"`
	Status       int64          `gorm:"default:1;comment:账户状态 1:正常 0:禁用" json:"status"`
	Avatar       string         `gorm:"size:255;comment:头像URL" json:"avatar"`
	Role         int64          `gorm:"not null;default:4;comment:身份角色 1:超级管理员 2:经理 3:主管 4:业务员（默认业务员）" json:"role"`
	MyInviteCode string         `gorm:"size:6;uniqueIndex:idx_admin_users_my_invite_code;comment:我的邀请码" json:"my_invite_code"`
	ParentID     *uint          `gorm:"index:idx_admin_users_parent_id;comment:上级用户ID" json:"parent_id"`
	CreatedAt    time.Time      `gorm:"type:datetime(3)" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"type:datetime(3)" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"type:datetime(3);index:idx_admin_users_deleted_at;comment:软删除时间" json:"-"`
}

// TableName 指定表名
func (AdminUser) TableName() string {
	return "admin_users"
}

// TableComment 表注释
func (AdminUser) TableComment() string {
	return "邀请码管理表 - 存储邀请码信息，用于用户注册时的邀请码校验"
}

// 管理员角色常量（使用int64枚举）
const (
	RoleSuperAdmin int64 = 1 // 超级管理员
	RoleManager    int64 = 2 // 经理
	RoleSupervisor int64 = 3 // 主管
	RoleSalesman   int64 = 4 // 业务员
)

// 角色名称映射
var RoleNames = map[int64]string{
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
func ValidateRoleID(roleID int64) bool {
	_, exists := RoleNames[roleID]
	return exists
}

// GetRoleIDByName 根据角色名称获取角色ID
func GetRoleIDByName(roleName string) (int64, bool) {
	for roleID, name := range RoleNames {
		if name == roleName {
			return roleID, true
		}
	}
	return 0, false
}
