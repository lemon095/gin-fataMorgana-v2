package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// Token过期时间
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	JWTSecret          []byte
)

// InitJWT 初始化JWT配置
func InitJWT(secret string, accessTokenExpire, refreshTokenExpire int) {
	JWTSecret = []byte(secret)
	AccessTokenExpiry = time.Duration(accessTokenExpire) * time.Second
	RefreshTokenExpiry = time.Duration(refreshTokenExpire) * time.Second
}

// Claims JWT声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Uid      string `json:"uid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateAccessToken 生成访问令牌
func GenerateAccessToken(userID uint, uid string, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Uid:      uid,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
			Issuer:    "gin-fataMorgana",
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken(userID uint, uid string, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Uid:      uid,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
			Issuer:    "gin-fataMorgana",
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// ParseToken 解析令牌
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, NewAppError(CodeTokenInvalid, "无效的令牌")
}

// ValidateToken 验证令牌
func ValidateToken(tokenString string) (*Claims, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 检查令牌是否过期
	if time.Now().UTC().Unix() > claims.ExpiresAt.Unix() {
		return nil, NewAppError(CodeTokenExpired, "令牌已过期")
	}

	return claims, nil
}
