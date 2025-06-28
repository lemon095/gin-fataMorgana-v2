package database

import (
	"context"
	"gin-fataMorgana/models"

	"gorm.io/gorm"
)

// AdminUserRepository 管理员用户仓库（仅用于邀请码校验）
type AdminUserRepository struct {
	db *gorm.DB
}

// NewAdminUserRepository 创建管理员用户仓库实例
func NewAdminUserRepository() *AdminUserRepository {
	return &AdminUserRepository{
		db: DB,
	}
}

// GetByInviteCode 根据邀请码获取管理员用户
func (r *AdminUserRepository) GetByInviteCode(ctx context.Context, inviteCode string) (*models.AdminUser, error) {
	var adminUser models.AdminUser
	err := r.db.WithContext(ctx).Where("my_invite_code = ?", inviteCode).First(&adminUser).Error
	if err != nil {
		return nil, err
	}
	return &adminUser, nil
}

// InviteCodeExists 检查邀请码是否存在
func (r *AdminUserRepository) InviteCodeExists(ctx context.Context, inviteCode string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.AdminUser{}).Where("my_invite_code = ?", inviteCode).Count(&count).Error
	return count > 0, err
}

// GetActiveInviteCode 获取活跃的邀请码（用于校验）
func (r *AdminUserRepository) GetActiveInviteCode(ctx context.Context, inviteCode string) (*models.AdminUser, error) {
	var adminUser models.AdminUser
	err := r.db.WithContext(ctx).
		Where("my_invite_code = ? AND status = ?", inviteCode, 1).
		First(&adminUser).Error
	if err != nil {
		return nil, err
	}
	return &adminUser, nil
}

// Create 创建管理员用户（用于初始化）
func (r *AdminUserRepository) Create(ctx context.Context, adminUser *models.AdminUser) error {
	return r.db.WithContext(ctx).Create(adminUser).Error
}

// GetByUsername 根据用户名获取管理员用户（用于登录校验）
func (r *AdminUserRepository) GetByUsername(ctx context.Context, username string) (*models.AdminUser, error) {
	var adminUser models.AdminUser
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&adminUser).Error
	if err != nil {
		return nil, err
	}
	return &adminUser, nil
}

// GetByID 根据ID获取管理员用户
func (r *AdminUserRepository) GetByID(ctx context.Context, id uint) (*models.AdminUser, error) {
	var adminUser models.AdminUser
	err := r.db.WithContext(ctx).First(&adminUser, id).Error
	if err != nil {
		return nil, err
	}
	return &adminUser, nil
}

// GetByAdminID 根据管理员ID获取管理员用户
func (r *AdminUserRepository) GetByAdminID(ctx context.Context, adminID string) (*models.AdminUser, error) {
	var adminUser models.AdminUser
	err := r.db.WithContext(ctx).Where("admin_id = ?", adminID).First(&adminUser).Error
	if err != nil {
		return nil, err
	}
	return &adminUser, nil
}

// Update 更新管理员用户
func (r *AdminUserRepository) Update(ctx context.Context, adminUser *models.AdminUser) error {
	return r.db.WithContext(ctx).Save(adminUser).Error
}

// Delete 软删除管理员用户
func (r *AdminUserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.AdminUser{}, id).Error
}

// List 获取管理员用户列表
func (r *AdminUserRepository) List(ctx context.Context, limit, offset int, role int) ([]models.AdminUser, error) {
	var adminUsers []models.AdminUser
	query := r.db.WithContext(ctx)
	
	if role > 0 {
		query = query.Where("role = ?", role)
	}
	
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&adminUsers).Error
	return adminUsers, err
}

// Count 统计管理员用户数量
func (r *AdminUserRepository) Count(ctx context.Context, role int) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.AdminUser{})
	
	if role > 0 {
		query = query.Where("role = ?", role)
	}
	
	err := query.Count(&count).Error
	return count, err
}

// Exists 检查记录是否存在
func (r *AdminUserRepository) Exists(ctx context.Context, conditions map[string]interface{}) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.AdminUser{}).Where(conditions).Count(&count).Error
	return count > 0, err
}

// AdminIDExists 检查管理员ID是否存在
func (r *AdminUserRepository) AdminIDExists(ctx context.Context, adminID string) (bool, error) {
	return r.Exists(ctx, map[string]interface{}{"admin_id": adminID})
}

// UsernameExists 检查用户名是否存在
func (r *AdminUserRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	return r.Exists(ctx, map[string]interface{}{"username": username})
}

// CheckUsernameAndInviteCode 批量检查用户名和邀请码是否存在（带缓存）
func (r *AdminUserRepository) CheckUsernameAndInviteCode(ctx context.Context, username, inviteCode string) (usernameExists bool, inviteCodeExists bool, err error) {
	// 先尝试从缓存获取用户名检查结果
	if cachedExists, hit, err := GetCachedUsernameExists(ctx, username); err == nil && hit {
		usernameExists = cachedExists
	} else {
		// 缓存未命中，查询数据库
		var count int64
		err = r.db.WithContext(ctx).Model(&models.AdminUser{}).Where("username = ?", username).Count(&count).Error
		if err != nil {
			return false, false, err
		}
		usernameExists = count > 0

		// 缓存结果
		CacheUsernameExists(ctx, username, usernameExists)
	}

	// 如果提供了邀请码，检查邀请码是否存在
	if inviteCode != "" {
		// 先尝试从缓存获取邀请码检查结果
		if cachedExists, hit, err := GetCachedInviteCodeExists(ctx, inviteCode); err == nil && hit {
			inviteCodeExists = cachedExists
		} else {
			// 缓存未命中，查询数据库
			var count int64
			err = r.db.WithContext(ctx).Model(&models.AdminUser{}).Where("my_invite_code = ?", inviteCode).Count(&count).Error
			if err != nil {
				return usernameExists, false, err
			}
			inviteCodeExists = count > 0

			// 缓存结果
			CacheInviteCodeExists(ctx, inviteCode, inviteCodeExists)
		}
	} else {
		inviteCodeExists = true // 如果没有提供邀请码，认为验证通过
	}

	return usernameExists, inviteCodeExists, nil
}

// BatchCheckUsernames 批量检查用户名是否存在（优化版本）
func (r *AdminUserRepository) BatchCheckUsernames(ctx context.Context, usernames []string) (map[string]bool, error) {
	if len(usernames) == 0 {
		return make(map[string]bool), nil
	}

	result := make(map[string]bool)
	uncachedUsernames := make([]string, 0)

	// 先检查缓存
	for _, username := range usernames {
		if cachedExists, hit, err := GetCachedUsernameExists(ctx, username); err == nil && hit {
			result[username] = cachedExists
		} else {
			uncachedUsernames = append(uncachedUsernames, username)
		}
	}

	// 批量查询未缓存的用户名
	if len(uncachedUsernames) > 0 {
		var adminUsers []models.AdminUser
		err := r.db.WithContext(ctx).Select("username").Where("username IN ?", uncachedUsernames).Find(&adminUsers).Error
		if err != nil {
			return nil, err
		}

		// 构建查询结果
		existingUsernames := make(map[string]bool)
		for _, adminUser := range adminUsers {
			existingUsernames[adminUser.Username] = true
		}

		// 更新结果并缓存
		for _, username := range uncachedUsernames {
			exists := existingUsernames[username]
			result[username] = exists
			CacheUsernameExists(ctx, username, exists)
		}
	}

	return result, nil
}

// BatchCheckInviteCodes 批量检查邀请码是否存在（优化版本）
func (r *AdminUserRepository) BatchCheckInviteCodes(ctx context.Context, inviteCodes []string) (map[string]bool, error) {
	if len(inviteCodes) == 0 {
		return make(map[string]bool), nil
	}

	result := make(map[string]bool)
	uncachedCodes := make([]string, 0)

	// 先检查缓存
	for _, code := range inviteCodes {
		if cachedExists, hit, err := GetCachedInviteCodeExists(ctx, code); err == nil && hit {
			result[code] = cachedExists
		} else {
			uncachedCodes = append(uncachedCodes, code)
		}
	}

	// 批量查询未缓存的邀请码
	if len(uncachedCodes) > 0 {
		var adminUsers []models.AdminUser
		err := r.db.WithContext(ctx).Select("my_invite_code").Where("my_invite_code IN ?", uncachedCodes).Find(&adminUsers).Error
		if err != nil {
			return nil, err
		}

		// 构建查询结果
		existingCodes := make(map[string]bool)
		for _, adminUser := range adminUsers {
			existingCodes[adminUser.MyInviteCode] = true
		}

		// 更新结果并缓存
		for _, code := range uncachedCodes {
			exists := existingCodes[code]
			result[code] = exists
			CacheInviteCodeExists(ctx, code, exists)
		}
	}

	return result, nil
}

// GetActiveAdmins 获取活跃管理员列表
func (r *AdminUserRepository) GetActiveAdmins(ctx context.Context, limit, offset int) ([]models.AdminUser, error) {
	var adminUsers []models.AdminUser
	err := r.db.WithContext(ctx).
		Where("status = ?", 1).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&adminUsers).Error
	return adminUsers, err
}

// GetAdminsByRole 根据角色获取管理员列表
func (r *AdminUserRepository) GetAdminsByRole(ctx context.Context, role int, limit, offset int) ([]models.AdminUser, error) {
	var adminUsers []models.AdminUser
	err := r.db.WithContext(ctx).
		Where("role = ?", role).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&adminUsers).Error
	return adminUsers, err
}

// GetAdminsByRoleLevel 根据角色等级获取管理员列表
func (r *AdminUserRepository) GetAdminsByRoleLevel(ctx context.Context, minLevel, maxLevel int, limit, offset int) ([]models.AdminUser, error) {
	var adminUsers []models.AdminUser
	
	// 直接使用角色ID范围查询
	err := r.db.WithContext(ctx).
		Where("role >= ? AND role <= ?", minLevel, maxLevel).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&adminUsers).Error
	return adminUsers, err
}
