package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"community-backward/handlers"
	"community-backward/middleware"
	"community-backward/utils"
)

// SetupRouter 配置路由
func SetupRouter(jwtUtil *utils.JWTUtil) *gin.Engine {
	// 创建Gin引擎
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 创建处理器
	authHandler := handlers.NewAuthHandler(jwtUtil)
	dashboardHandler := handlers.NewDashboardHandler()
	residentHandler := handlers.NewResidentHandler()

	// API路由组
	api := r.Group("/api")
	{
		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		// 受保护的路由
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(jwtUtil))
		{
			// 仪表盘相关路由
			protected.GET("/dashboard", dashboardHandler.GetDashboard)
			protected.GET("/dashboard/income", dashboardHandler.GetIncomeData)

			// 住户相关路由
			residents := protected.Group("/residents")
			{
				residents.GET("", residentHandler.GetResidents)
				residents.GET("/:id", residentHandler.GetResidentByID)
				residents.POST("", residentHandler.CreateResident)
				residents.PUT("/:id", residentHandler.UpdateResident)
				residents.DELETE("/:id", residentHandler.DeleteResident)
			}
		}
	}

	return r
}
