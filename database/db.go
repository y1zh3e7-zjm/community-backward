package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	// DB 是全局数据库连接
	DB *sql.DB
)

// InitDB 初始化数据库连接
func InitDB() error {
	// 从环境变量获取数据库配置
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "community")

	// 打印数据库配置信息
	log.Printf("数据库配置: 用户=%s, 主机=%s, 端口=%s, 数据库=%s", dbUser, dbHost, dbPort, dbName)

	// 构建DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// 打开数据库连接
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("无法连接到数据库: %v", err)
	}

	// 设置连接池参数
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// 测试连接
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	// 检查数据库表是否存在
	checkTables()

	log.Println("数据库连接成功")
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("数据库连接已关闭")
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// checkTables 检查数据库表是否存在
func checkTables() {
	// 检查residents表
	var residentCount int
	err := DB.QueryRow("SELECT COUNT(*) FROM residents").Scan(&residentCount)
	if err != nil {
		log.Printf("检查residents表失败: %v", err)
	} else {
		log.Printf("residents表中有 %d 条记录", residentCount)
	}

	// 检查property_fees表
	var propertyFeeCount int
	err = DB.QueryRow("SELECT COUNT(*) FROM property_fees").Scan(&propertyFeeCount)
	if err != nil {
		log.Printf("检查property_fees表失败: %v", err)
	} else {
		log.Printf("property_fees表中有 %d 条记录", propertyFeeCount)
	}

	// 检查property_fee_monthly_stats表
	var statsCount int
	err = DB.QueryRow("SELECT COUNT(*) FROM property_fee_monthly_stats").Scan(&statsCount)
	if err != nil {
		log.Printf("检查property_fee_monthly_stats表失败: %v", err)
	} else {
		log.Printf("property_fee_monthly_stats表中有 %d 条记录", statsCount)
	}
}
