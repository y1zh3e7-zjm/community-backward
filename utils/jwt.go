package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTUtil 提供JWT相关功能
type JWTUtil struct {
	Secret []byte
}

// NewJWTUtil 创建新的JWTUtil实例
func NewJWTUtil(secret []byte) *JWTUtil {
	return &JWTUtil{
		Secret: secret,
	}
}

// GenerateToken 生成JWT令牌
func (j *JWTUtil) GenerateToken(userID uint, username string, remember bool) (string, error) {
	// 设置过期时间
	expirationTime := time.Now().Add(24 * time.Hour)
	if remember {
		expirationTime = time.Now().Add(7 * 24 * time.Hour) // 记住密码时token有效期为7天
	}

	// 创建claims
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      expirationTime.Unix(),
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// 签名token
	return token.SignedString(j.Secret)
}

// ValidateToken 验证JWT令牌
func (j *JWTUtil) ValidateToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	// 解析token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.Secret, nil
	})

	if err != nil || !token.Valid {
		return nil, nil, err
	}

	// 提取claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, jwt.ErrTokenInvalidClaims
	}

	return token, claims, nil
}
