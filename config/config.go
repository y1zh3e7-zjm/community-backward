package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config 存储应用程序配置
type Config struct {
	Port      string
	JWTSecret []byte
	GinMode   string
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 加载.env文件
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default environment")
	}

	// 获取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// 获取JWT密钥
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your_jwt_secret_key" // 默认密钥，生产环境应该更改
	}

	// 获取Gin模式
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "release"
	}

	return &Config{
		Port:      port,
		JWTSecret: []byte(jwtSecret),
		GinMode:   ginMode,
	}
}
