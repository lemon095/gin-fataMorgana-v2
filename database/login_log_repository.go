package database

import (
	"context"
	"gin-fataMorgana/models"
	"time"
)

// LoginLogRepository 登录记录仓库
type LoginLogRepository struct {
	*BaseRepository
}

// NewLoginLogRepository 创建登录记录仓库
func NewLoginLogRepository() *LoginLogRepository {
	return &LoginLogRepository{
		BaseRepository: NewBaseRepository(),
	}
}

// Create 创建登录记录
func (r *LoginLogRepository) Create(ctx context.Context, log *models.UserLoginLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// GetUserLoginHistory 获取用户登录历史
func (r *LoginLogRepository) GetUserLoginHistory(ctx context.Context, uid string, limit, offset int) ([]models.UserLoginLog, error) {
	var logs []models.UserLoginLog
	err := r.db.WithContext(ctx).
		Where("uid = ?", uid).
		Order("login_time DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, err
}

// GetUserLastLogin 获取用户最后登录记录
func (r *LoginLogRepository) GetUserLastLogin(ctx context.Context, uid string) (*models.UserLoginLog, error) {
	var log models.UserLoginLog
	err := r.db.WithContext(ctx).
		Where("uid = ? AND status = ?", uid, 1).
		Order("login_time DESC").
		First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetLoginStats 获取登录统计信息
func (r *LoginLogRepository) GetLoginStats(ctx context.Context, uid string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总登录次数
	var totalLogins int64
	err := r.db.WithContext(ctx).Model(&models.UserLoginLog{}).
		Where("uid = ?", uid).Count(&totalLogins).Error
	if err != nil {
		return nil, err
	}
	stats["total_logins"] = totalLogins

	// 成功登录次数
	var successLogins int64
	err = r.db.WithContext(ctx).Model(&models.UserLoginLog{}).
		Where("uid = ? AND status = ?", uid, 1).Count(&successLogins).Error
	if err != nil {
		return nil, err
	}
	stats["success_logins"] = successLogins

	// 失败登录次数
	stats["failed_logins"] = totalLogins - successLogins

	// 最后登录时间
	var lastLogin models.UserLoginLog
	err = r.db.WithContext(ctx).
		Where("uid = ? AND status = ?", uid, 1).
		Order("login_time DESC").
		First(&lastLogin).Error
	if err == nil {
		stats["last_login_time"] = lastLogin.LoginTime
		stats["last_login_ip"] = lastLogin.LoginIP
	}

	// 最近7天登录次数
	var recentLogins int64
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	err = r.db.WithContext(ctx).Model(&models.UserLoginLog{}).
		Where("uid = ? AND login_time >= ?", uid, sevenDaysAgo).Count(&recentLogins).Error
	if err != nil {
		return nil, err
	}
	stats["recent_7_days_logins"] = recentLogins

	return stats, nil
}

// GetLoginLogsByTimeRange 按时间范围获取登录记录
func (r *LoginLogRepository) GetLoginLogsByTimeRange(ctx context.Context, uid string, startTime, endTime time.Time) ([]models.UserLoginLog, error) {
	var logs []models.UserLoginLog
	err := r.db.WithContext(ctx).
		Where("uid = ? AND login_time BETWEEN ? AND ?", uid, startTime, endTime).
		Order("login_time DESC").
		Find(&logs).Error
	return logs, err
}

// GetLoginLogsByIP 按IP地址获取登录记录
func (r *LoginLogRepository) GetLoginLogsByIP(ctx context.Context, uid, ip string) ([]models.UserLoginLog, error) {
	var logs []models.UserLoginLog
	err := r.db.WithContext(ctx).
		Where("uid = ? AND login_ip = ?", uid, ip).
		Order("login_time DESC").
		Find(&logs).Error
	return logs, err
}

// GetFailedLoginAttempts 获取失败登录尝试次数
func (r *LoginLogRepository) GetFailedLoginAttempts(ctx context.Context, uid string, since time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.UserLoginLog{}).
		Where("uid = ? AND status = ? AND login_time >= ?", uid, 0, since).
		Count(&count).Error
	return count, err
}

// CleanOldLogs 清理旧登录记录
func (r *LoginLogRepository) CleanOldLogs(ctx context.Context, before time.Time) error {
	return r.db.WithContext(ctx).
		Where("login_time < ?", before).
		Delete(&models.UserLoginLog{}).Error
}
