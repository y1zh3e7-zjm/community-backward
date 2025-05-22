package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"community-backward/models"
)

// ResidentHandler 处理住户相关请求
type ResidentHandler struct{}

// NewResidentHandler 创建新的ResidentHandler
func NewResidentHandler() *ResidentHandler {
	return &ResidentHandler{}
}

// GetResidents 获取住户列表
func (h *ResidentHandler) GetResidents(c *gin.Context) {
	fmt.Println("接收到获取住户列表请求")

	// 打印所有请求参数
	queryParams := c.Request.URL.Query()
	fmt.Println("请求参数:", queryParams)

	// 获取分页参数
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	fmt.Println("页码:", page)

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	fmt.Println("每页条数:", pageSize)

	// 获取过滤参数
	filters := make(map[string]interface{})

	// 姓名过滤
	if name := c.Query("name"); name != "" {
		filters["name"] = name
		fmt.Println("姓名过滤:", name)
	}

	// 楼栋过滤
	if building := c.Query("building"); building != "" {
		filters["building"] = building
		fmt.Println("楼栋过滤:", building)
	}

	// 单元过滤
	if unit := c.Query("unit"); unit != "" {
		filters["unit"] = unit
		fmt.Println("单元过滤:", unit)
	}

	// 面积过滤
	if area := c.Query("area"); area != "" {
		if areaFloat, err := strconv.ParseFloat(area, 64); err == nil {
			filters["area"] = areaFloat
			fmt.Println("最小面积过滤:", areaFloat)
		}
	}

	if areaMax := c.Query("areaMax"); areaMax != "" {
		if areaMaxFloat, err := strconv.ParseFloat(areaMax, 64); err == nil {
			filters["areaMax"] = areaMaxFloat
			fmt.Println("最大面积过滤:", areaMaxFloat)
		}
	}

	// 排序参数
	if sortField := c.Query("sortField"); sortField != "" {
		filters["sortField"] = sortField
		fmt.Println("排序字段:", sortField)
	}

	if sortOrder := c.Query("sortOrder"); sortOrder != "" {
		filters["sortOrder"] = sortOrder
		fmt.Println("排序方向:", sortOrder)
	}

	fmt.Println("最终过滤条件:", filters)

	// 获取住户列表
	residents, total, err := models.GetResidents(page, pageSize, filters)
	if err != nil {
		fmt.Println("获取住户列表失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取住户列表失败: " + err.Error(),
		})
		return
	}

	fmt.Println("查询成功，返回住户数量:", len(residents))
	fmt.Println("总数:", total)

	// 构建响应
	response := models.ResidentListResponse{
		Success: true,
	}
	response.Data.List = residents
	response.Data.Total = total

	c.JSON(http.StatusOK, response)
}

// GetResidentByID 获取住户详情
func (h *ResidentHandler) GetResidentByID(c *gin.Context) {
	// 获取ID参数
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的住户ID",
		})
		return
	}

	// 获取住户
	resident, err := models.GetResidentByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取住户详情失败: " + err.Error(),
		})
		return
	}

	if resident == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "住户不存在",
		})
		return
	}

	// 构建响应
	c.JSON(http.StatusOK, models.ResidentResponse{
		Success: true,
		Data:    *resident,
	})
}

// CreateResident 创建住户
func (h *ResidentHandler) CreateResident(c *gin.Context) {
	fmt.Println("接收到创建住户请求")

	// 打印请求体
	body, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	fmt.Println("请求体:", string(body))

	var req models.ResidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("解析请求参数失败:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的请求参数: " + err.Error(),
		})
		return
	}

	fmt.Println("解析后的请求参数:", req)

	// 创建住户
	id, err := models.CreateResident(req)
	if err != nil {
		fmt.Println("创建住户失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建住户失败: " + err.Error(),
		})
		return
	}

	fmt.Println("住户创建成功，ID:", id)

	// 获取新创建的住户
	resident, err := models.GetResidentByID(id)
	if err != nil {
		fmt.Println("获取新创建的住户失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取新创建的住户失败: " + err.Error(),
		})
		return
	}

	fmt.Println("获取到的新住户:", resident)

	// 构建响应
	response := models.ResidentResponse{
		Success: true,
		Message: "住户创建成功",
		Data:    *resident,
	}
	fmt.Println("返回响应:", response)

	c.JSON(http.StatusCreated, response)
}

// UpdateResident 更新住户
func (h *ResidentHandler) UpdateResident(c *gin.Context) {
	// 获取ID参数
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的住户ID",
		})
		return
	}

	// 检查住户是否存在
	resident, err := models.GetResidentByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取住户详情失败: " + err.Error(),
		})
		return
	}

	if resident == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "住户不存在",
		})
		return
	}

	// 解析请求体
	var req models.ResidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的请求参数: " + err.Error(),
		})
		return
	}

	// 更新住户
	if err := models.UpdateResident(id, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新住户失败: " + err.Error(),
		})
		return
	}

	// 获取更新后的住户
	updatedResident, err := models.GetResidentByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取更新后的住户失败: " + err.Error(),
		})
		return
	}

	// 构建响应
	c.JSON(http.StatusOK, models.ResidentResponse{
		Success: true,
		Message: "住户更新成功",
		Data:    *updatedResident,
	})
}

// DeleteResident 删除住户
func (h *ResidentHandler) DeleteResident(c *gin.Context) {
	// 获取ID参数
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的住户ID",
		})
		return
	}

	// 检查住户是否存在
	resident, err := models.GetResidentByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取住户详情失败: " + err.Error(),
		})
		return
	}

	if resident == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "住户不存在",
		})
		return
	}

	// 删除住户
	if err := models.DeleteResident(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除住户失败: " + err.Error(),
		})
		return
	}

	// 构建响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "住户删除成功",
	})
}
