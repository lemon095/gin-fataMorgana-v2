package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gin-fataMorgana/config"
	"gin-fataMorgana/database"
	"gin-fataMorgana/utils"
)

// TokenService Token管理服务
type TokenService struct{}

// NewTokenService 创建Token服务实例
func NewTokenService() *TokenService {
	return &TokenService{}
}

// TokenInfo Token信息
type TokenInfo struct {
	TokenHash   string    `json:"token_hash"`
	LoginTime   time.Time `json:"login_time"`
	DeviceInfo  string    `json:"device_info"`
	LoginIP     string    `json:"login_ip"`
	UserAgent   string    `json:"user_agent"`
}

// SessionInfo 会话信息
type SessionInfo struct {
	CurrentTokenHash string    `json:"current_token_hash"`
	LastLoginTime    time.Time `json:"last_login_time"`
	DeviceInfo       string    `json:"device_info"`
	LoginIP          string    `json:"login_ip"`
	UserAgent        string    `json:"user_agent"`
}

// 生成token哈希
func (s *TokenService) generateTokenHash(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// 获取用户活跃token的Redis key
func (s *TokenService) getUserActiveTokenKey(uid string) string {
	return fmt.Sprintf("user:active_token:%s", uid)
}

// 获取token黑名单的Redis key
func (s *TokenService) getTokenBlacklistKey(tokenHash string) string {
	return fmt.Sprintf("token:blacklist:%s", tokenHash)
}

// 获取用户会话信息的Redis key
func (s *TokenService) getUserSessionKey(uid string) string {
	return fmt.Sprintf("user:session:%s", uid)
}

// GetUserActiveToken 获取用户当前活跃token
func (s *TokenService) GetUserActiveToken(ctx context.Context, uid string) (*TokenInfo, error) {
	key := s.getUserActiveTokenKey(uid)
	
	value, err := database.GetKey(ctx, key)
	if err != nil {
		return nil, err
	}
	
	if value == "" {
		return nil, nil
	}
	
	var tokenInfo TokenInfo
	if err := json.Unmarshal([]byte(value), &tokenInfo); err != nil {
		return nil, err
	}
	
	return &tokenInfo, nil
}

// SetUserActiveToken 设置用户活跃token
func (s *TokenService) SetUserActiveToken(ctx context.Context, uid, token, deviceInfo, loginIP, userAgent string) error {
	tokenHash := s.generateTokenHash(token)
	
	tokenInfo := &TokenInfo{
		TokenHash:  tokenHash,
		LoginTime:  time.Now().UTC(),
		DeviceInfo: deviceInfo,
		LoginIP:    loginIP,
		UserAgent:  userAgent,
	}
	
	// 序列化为JSON
	data, err := json.Marshal(tokenInfo)
	if err != nil {
		return err
	}
	
	// 设置活跃token，过期时间与token一致
	expiration := time.Duration(config.GlobalConfig.JWT.AccessTokenExpire) * time.Second
	key := s.getUserActiveTokenKey(uid)
	
	return database.SetKey(ctx, key, string(data), expiration)
}

// AddTokenToBlacklist 将token加入黑名单
func (s *TokenService) AddTokenToBlacklist(ctx context.Context, token string) error {
	tokenHash := s.generateTokenHash(token)
	key := s.getTokenBlacklistKey(tokenHash)
	
	// 黑名单过期时间比token过期时间稍长，确保覆盖
	expiration := time.Duration(config.GlobalConfig.JWT.AccessTokenExpire+300) * time.Second
	
	return database.SetKey(ctx, key, time.Now().UTC().Unix(), expiration)
}

// IsTokenBlacklisted 检查token是否在黑名单中
func (s *TokenService) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	tokenHash := s.generateTokenHash(token)
	key := s.getTokenBlacklistKey(tokenHash)
	
	exists, err := database.ExistsKey(ctx, key)
	if err != nil {
		return false, err
	}
	
	return exists, nil
}

// IsActiveToken 检查token是否为当前活跃token
func (s *TokenService) IsActiveToken(ctx context.Context, uid, token string) (bool, error) {
	activeToken, err := s.GetUserActiveToken(ctx, uid)
	if err != nil {
		return false, err
	}
	
	if activeToken == nil {
		return false, nil
	}
	
	tokenHash := s.generateTokenHash(token)
	return activeToken.TokenHash == tokenHash, nil
}

// RevokeUserSession 撤销用户会话
func (s *TokenService) RevokeUserSession(ctx context.Context, uid string) error {
	// 获取当前活跃token
	activeToken, err := s.GetUserActiveToken(ctx, uid)
	if err != nil {
		return err
	}
	
	if activeToken != nil {
		// 将活跃token加入黑名单
		key := s.getTokenBlacklistKey(activeToken.TokenHash)
		expiration := time.Duration(config.GlobalConfig.JWT.AccessTokenExpire+300) * time.Second
		database.SetKey(ctx, key, time.Now().UTC().Unix(), expiration)
	}
	
	// 删除活跃token记录
	activeKey := s.getUserActiveTokenKey(uid)
	database.DelKey(ctx, activeKey)
	
	return nil
}

// GetUserSessionInfo 获取用户会话信息
func (s *TokenService) GetUserSessionInfo(ctx context.Context, uid string) (*SessionInfo, error) {
	activeToken, err := s.GetUserActiveToken(ctx, uid)
	if err != nil {
		return nil, err
	}
	
	if activeToken == nil {
		return nil, nil
	}
	
	return &SessionInfo{
		CurrentTokenHash: activeToken.TokenHash,
		LastLoginTime:    activeToken.LoginTime,
		DeviceInfo:       activeToken.DeviceInfo,
		LoginIP:          activeToken.LoginIP,
		UserAgent:        activeToken.UserAgent,
	}, nil
}

// CleanupExpiredSessions 清理过期会话
func (s *TokenService) CleanupExpiredSessions(ctx context.Context) error {
	// 这里可以实现定期清理逻辑
	// 由于Redis会自动过期，主要清理逻辑在token验证时进行
	return nil
}

// ValidateTokenWithBlacklist 验证令牌（包含黑名单检查）
func (s *TokenService) ValidateTokenWithBlacklist(ctx context.Context, tokenString string) (*utils.Claims, error) {
	// 首先进行基本的JWT验证
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 检查token是否在黑名单中
	isBlacklisted, err := s.IsTokenBlacklisted(ctx, tokenString)
	if err != nil {
		return nil, errors.New("token验证失败")
	}

	if isBlacklisted {
		return nil, errors.New("您的账号已在其他设备登录，请重新登录")
	}

	// 检查是否为当前活跃token
	isActive, err := s.IsActiveToken(ctx, claims.Uid, tokenString)
	if err != nil {
		return nil, errors.New("token验证失败")
	}

	if !isActive {
		return nil, errors.New("您的账号已在其他设备登录，请重新登录")
	}

	return claims, nil
} 