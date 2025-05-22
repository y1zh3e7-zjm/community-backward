package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// HashPassword 使用MD5对密码进行哈希
func HashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

// CheckPasswordHash 检查密码是否与哈希匹配
func CheckPasswordHash(password, hash string) bool {
	return HashPassword(password) == hash
}
