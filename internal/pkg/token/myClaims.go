package token

import (
	"time"
	"wallpaper/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	jwt.RegisteredClaims
	UserID     uint   `json:"user_id"`
	Name       string `json:"name"`
	Unionid    string `json:"unionid"`
	SessionKey string `json:"session_key"`
}

// GenerateToken 生成访问令牌，返回 token 字符串和过期时间戳
func GenerateToken(config *config.Auth, claims MyClaims) (string, int64, error) {
	expireAt := time.Now().Add(time.Second * time.Duration(config.AccessExpire))
	if claims.RegisteredClaims.ExpiresAt == nil {
		claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(expireAt)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.AccessSecret))
	if err != nil {
		return "", 0, err
	}
	return tokenStr, expireAt.Unix(), nil
}

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken(config *config.Auth, userID uint) (string, error) {
	expireAt := time.Now().Add(time.Second * time.Duration(config.RefreshExpire))
	claims := MyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.RefreshSecret))
}

func ParseToken(secret string, tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// ParseRefreshToken 解析刷新令牌
func ParseRefreshToken(secret string, tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}
