package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"community-backward/config"
	"community-backward/database"
	"community-backward/routes"
	"community-backward/utils"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 设置Gin模式
	gin.SetMode(cfg.GinMode)

	// 初始化数据库连接
	if err := database.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer database.CloseDB()

	// 创建JWT工具
	jwtUtil := utils.NewJWTUtil(cfg.JWTSecret)

	// 设置路由
	r := routes.SetupRouter(jwtUtil)

	// 启动服务器
	log.Printf("Server running on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

