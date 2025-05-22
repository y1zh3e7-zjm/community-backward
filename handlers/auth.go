package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"community-backward/models"
	"community-backward/utils"
)

// AuthHandler 处理认证相关请求
type AuthHandler struct {
	JWTUtil *utils.JWTUtil
}

// NewAuthHandler 创建新的AuthHandler
func NewAuthHandler(jwtUtil *utils.JWTUtil) *AuthHandler {
	return &AuthHandler{
		JWTUtil: jwtUtil,
	}
}

// Login 处理登录请求
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的请求参数",
		})
		return
	}

	// 验证用户
	user, err := models.FindUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "服务器错误",
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "账号或密码错误",
		})
		return
	}

	// 生成JWT
	tokenString, err := h.JWTUtil.GenerateToken(user.ID, user.Username, req.Remember)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "生成token失败",
		})
		return
	}

	// 返回token
	c.JSON(http.StatusOK, models.LoginResponse{
		Success: true,
		Token:   tokenString,
		Message: "登录成功",
	})
}
