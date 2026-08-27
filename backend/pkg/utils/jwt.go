package utils

import (
	"art-design-pro-api/internal/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint   `json:"userId"`
	UserName string `json:"userName"`
	jwt.RegisteredClaims
}

// GenerateToken 生成访问令牌
func GenerateToken(userID uint, userName string) (string, error) {
	return generate(userID, userName, time.Duration(config.AppConfig.JWT.ExpireHours)*time.Hour)
}

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken(userID uint, userName string) (string, error) {
	return generate(userID, userName, time.Duration(config.AppConfig.JWT.RefreshExpireHours)*time.Hour)
}

func generate(userID uint, userName string, expire time.Duration) (string, error) {
	claims := Claims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

// ParseToken 解析令牌
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
