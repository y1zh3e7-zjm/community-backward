package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"community-backward/utils"
)

// AuthMiddleware 创建认证中间件
func AuthMiddleware(jwtUtil *utils.JWTUtil) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "未提供认证信息",
			})
			return
		}

		// 提取Bearer token
		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "认证格式错误",
			})
			return
		}

		// 验证token
		_, claims, err := jwtUtil.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "无效的认证信息",
			})
			return
		}

		// 将用户信息存储在上下文中
		c.Set("user_id", claims["user_id"])
		c.Set("username", claims["username"])

		c.Next()
	}
}
