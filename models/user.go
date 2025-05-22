package models

import (
	"database/sql"
	"time"
	"community-backward/database"
	"community-backward/utils"
)

// User 表示用户模型
type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // 不在JSON中返回密码
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// LoginRequest 表示登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Remember bool   `json:"remember"`
}

// LoginResponse 表示登录响应
type LoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Message string `json:"message,omitempty"`
}

// FindUserByUsername 通过用户名查找用户
func FindUserByUsername(username string) (*User, error) {
	query := `SELECT id, username, password, email, created_at, updated_at FROM users WHERE username = ?`

	var user User
	err := database.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 用户不存在
		}
		return nil, err
	}

	return &user, nil
}

// FindUser 通过用户名和密码查找用户（用于登录验证）
func FindUser(username, password string) (*User, error) {
	user, err := FindUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	// 验证密码（MD5哈希比较）
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, nil // 密码不匹配
	}

	return user, nil
}

// GetUserByID 通过ID获取用户
func GetUserByID(id uint) (*User, error) {
	query := `SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?`

	var user User
	err := database.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 用户不存在
		}
		return nil, err
	}

	return &user, nil
}


