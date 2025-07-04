package services

import (
	"context"
	"gin-fataMorgana/database"
	"gin-fataMorgana/models"
	"time"
)

// LoginLogService 登录记录服务
type LoginLogService struct {
	loginLogRepo *database.LoginLogRepository
}

// NewLoginLogService 创建登录记录服务实例
func NewLoginLogService() *LoginLogService {
	return &LoginLogService{
		loginLogRepo: database.NewLoginLogRepository(),
	}
}

// RecordLogin 记录登录
func (s *LoginLogService) RecordLogin(ctx context.Context, user *models.User, loginIP, userAgent string, status int, failReason string) error {
	// 获取设备信息
	deviceInfo := s.extractDeviceInfo(userAgent)

	// 获取地理位置（这里可以集成IP地理位置服务）
	location := s.getLocationByIP(loginIP)

	log := &models.UserLoginLog{
		Uid:        user.Uid,
		Username:   user.Username,
		Email:      user.Email,
		LoginIP:    loginIP,
		UserAgent:  userAgent,
		LoginTime:  time.Now().UTC(),
		Status:     status,
		FailReason: failReason,
		DeviceInfo: deviceInfo,
		Location:   location,
	}

	return s.loginLogRepo.Create(ctx, log)
}

// GetUserLoginHistory 获取用户登录历史
func (s *LoginLogService) GetUserLoginHistory(ctx context.Context, uid string, page, size int) ([]models.UserLoginLogResponse, int64, error) {
	offset := (page - 1) * size

	logs, err := s.loginLogRepo.GetUserLoginHistory(ctx, uid, size, offset)
	if err != nil {
		return nil, 0, err
	}

	// 获取总数
	total, err := s.getUserLoginCount(ctx, uid)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	responses := make([]models.UserLoginLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = log.ToResponse()
	}

	return responses, total, nil
}

// GetUserLastLogin 获取用户最后登录信息
func (s *LoginLogService) GetUserLastLogin(ctx context.Context, uid string) (*models.UserLoginLogResponse, error) {
	log, err := s.loginLogRepo.GetUserLastLogin(ctx, uid)
	if err != nil {
		return nil, err
	}

	response := log.ToResponse()
	return &response, nil
}

// GetLoginStats 获取登录统计信息
func (s *LoginLogService) GetLoginStats(ctx context.Context, uid string) (map[string]interface{}, error) {
	return s.loginLogRepo.GetLoginStats(ctx, uid)
}

// GetLoginLogsByTimeRange 按时间范围获取登录记录
func (s *LoginLogService) GetLoginLogsByTimeRange(ctx context.Context, uid string, startTime, endTime time.Time) ([]models.UserLoginLogResponse, error) {
	logs, err := s.loginLogRepo.GetLoginLogsByTimeRange(ctx, uid, startTime, endTime)
	if err != nil {
		return nil, err
	}

	responses := make([]models.UserLoginLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = log.ToResponse()
	}

	return responses, nil
}

// GetLoginLogsByIP 按IP地址获取登录记录
func (s *LoginLogService) GetLoginLogsByIP(ctx context.Context, uid, ip string) ([]models.UserLoginLogResponse, error) {
	logs, err := s.loginLogRepo.GetLoginLogsByIP(ctx, uid, ip)
	if err != nil {
		return nil, err
	}

	responses := make([]models.UserLoginLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = log.ToResponse()
	}

	return responses, nil
}

// CheckFailedLoginAttempts 检查失败登录尝试次数
func (s *LoginLogService) CheckFailedLoginAttempts(ctx context.Context, uid string, maxAttempts int, lockoutDuration time.Duration) (bool, error) {
	since := time.Now().UTC().Add(-lockoutDuration)
	count, err := s.loginLogRepo.GetFailedLoginAttempts(ctx, uid, since)
	if err != nil {
		return false, err
	}

	return count >= int64(maxAttempts), nil
}

// CleanOldLogs 清理旧登录记录
func (s *LoginLogService) CleanOldLogs(ctx context.Context, days int) error {
	before := time.Now().UTC().AddDate(0, 0, -days)
	return s.loginLogRepo.CleanOldLogs(ctx, before)
}

// 辅助方法

// getUserLoginCount 获取用户登录记录总数
func (s *LoginLogService) getUserLoginCount(ctx context.Context, uid string) (int64, error) {
	// 这里可以添加缓存优化
	count, err := s.loginLogRepo.Count(ctx, map[string]interface{}{"uid": uid})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// extractDeviceInfo 提取设备信息
func (s *LoginLogService) extractDeviceInfo(userAgent string) string {
	// 简单的设备信息提取，可以集成更复杂的解析库
	if len(userAgent) > 200 {
		userAgent = userAgent[:200]
	}
	return userAgent
}

// getLocationByIP 根据IP获取地理位置
func (s *LoginLogService) getLocationByIP(ip string) string {
	// 这里可以集成IP地理位置服务，如GeoIP2、IP2Location等
	// 暂时返回空字符串，后续可以扩展
	return ""
}
