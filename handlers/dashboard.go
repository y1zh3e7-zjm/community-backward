package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"community-backward/models"
)

// DashboardHandler 处理仪表盘相关请求
type DashboardHandler struct{}

// NewDashboardHandler 创建新的DashboardHandler
func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

// GetDashboard 获取仪表盘数据
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	fmt.Println("接收到获取仪表盘数据请求")

	username, _ := c.Get("username")
	fmt.Printf("当前用户: %v\n", username)

	// 获取仪表盘统计数据
	stats, err := models.GetDashboardStats()
	if err != nil {
		fmt.Println("获取仪表盘统计数据失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取仪表盘数据失败: " + err.Error(),
		})
		return
	}

	fmt.Printf("返回仪表盘统计数据: %+v\n", stats)

	response := gin.H{
		"success": true,
		"data": gin.H{
			"message": "欢迎访问仪表盘",
			"username": username,
			"stats": stats,
		},
	}

	fmt.Printf("返回响应: %+v\n", response)
	c.JSON(http.StatusOK, response)
}

// GetIncomeData 获取物业费收入趋势数据
func (h *DashboardHandler) GetIncomeData(c *gin.Context) {
	fmt.Println("接收到获取物业费收入趋势数据请求")

	// 获取年份参数，默认为当前年份
	yearStr := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))
	fmt.Printf("请求参数 year: %s\n", yearStr)

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		fmt.Printf("无效的年份参数: %s, 错误: %v\n", yearStr, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的年份参数",
		})
		return
	}

	fmt.Printf("解析后的年份: %d\n", year)

	// 获取物业费收入趋势数据
	incomeData, err := models.GetIncomeData(year)
	if err != nil {
		fmt.Println("获取物业费收入趋势数据失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取物业费收入趋势数据失败: " + err.Error(),
		})
		return
	}

	fmt.Printf("返回物业费收入趋势数据: %+v\n", incomeData)

	response := gin.H{
		"success": true,
		"data": incomeData,
	}

	fmt.Printf("返回响应: %+v\n", response)
	c.JSON(http.StatusOK, response)
}
