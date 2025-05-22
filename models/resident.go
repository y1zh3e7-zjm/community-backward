package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"community-backward/database"
)

// Resident 表示住户模型
type Resident struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	BlockNumber string    `json:"building"` // 前端使用building字段
	UnitNumber  string    `json:"unit"`     // 前端使用unit字段
	HouseNumber string    `json:"room"`     // 前端使用room字段
	HouseArea   float64   `json:"area"`     // 前端使用area字段
	FareSum     float64   `json:"fee"`      // 前端使用fee字段
	Status      int       `json:"-"`        // 数据库中的状态值
	StatusText  string    `json:"status"`   // 前端显示的状态文本
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// ResidentRequest 表示住户请求
type ResidentRequest struct {
	Name        string  `json:"name" binding:"required"`
	BlockNumber string  `json:"building" binding:"required"` // 前端使用building字段
	UnitNumber  string  `json:"unit" binding:"required"`     // 前端使用unit字段
	HouseNumber string  `json:"room" binding:"required"`     // 前端使用room字段
	HouseArea   float64 `json:"area" binding:"required"`     // 前端使用area字段
	FareSum     float64 `json:"fee"`                         // 前端使用fee字段
	Status      int     `json:"status"`                      // 0表示欠费，1表示正常
}

// ResidentListResponse 表示住户列表响应
type ResidentListResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    struct {
		List  []Resident `json:"list"`
		Total int        `json:"total"`
	} `json:"data"`
}

// ResidentResponse 表示住户响应
type ResidentResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Data    Resident `json:"data,omitempty"`
}

// GetResidents 获取住户列表
func GetResidents(page, pageSize int, filters map[string]interface{}) ([]Resident, int, error) {
	// 打印接收到的筛选条件
	fmt.Println("接收到的筛选条件:", filters)

	// 构建查询
	query := `SELECT id, name, block_number, unit_number, house_number, house_area, fare_sum, status, created_at, updated_at FROM residents WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM residents WHERE 1=1`

	// 参数列表
	var args []interface{}
	var countArgs []interface{}

	// 添加过滤条件
	if name, ok := filters["name"].(string); ok && name != "" {
		query += ` AND name LIKE ?`
		countQuery += ` AND name LIKE ?`
		args = append(args, "%"+name+"%")
		countArgs = append(countArgs, "%"+name+"%")
		fmt.Println("添加姓名过滤条件:", name)
	}

	if building, ok := filters["building"].(string); ok && building != "" {
		// 检查building是否已经包含"栋"字
		buildingStr := building
		if !strings.HasSuffix(buildingStr, "栋") {
			buildingStr = buildingStr + "栋"
		}
		query += ` AND block_number = ?`
		countQuery += ` AND block_number = ?`
		args = append(args, buildingStr)
		countArgs = append(countArgs, buildingStr)
		fmt.Println("添加楼栋过滤条件:", buildingStr)
	}

	if unit, ok := filters["unit"].(string); ok && unit != "" {
		// 检查unit是否已经包含"单元"字
		unitStr := unit
		if !strings.HasSuffix(unitStr, "单元") {
			unitStr = unitStr + "单元"
		}
		query += ` AND unit_number = ?`
		countQuery += ` AND unit_number = ?`
		args = append(args, unitStr)
		countArgs = append(countArgs, unitStr)
		fmt.Println("添加单元过滤条件:", unitStr)
	}

	if area, ok := filters["area"].(float64); ok && area > 0 {
		query += ` AND house_area >= ?`
		countQuery += ` AND house_area >= ?`
		args = append(args, area)
		countArgs = append(countArgs, area)
	}

	if areaMax, ok := filters["areaMax"].(float64); ok && areaMax > 0 {
		query += ` AND house_area <= ?`
		countQuery += ` AND house_area <= ?`
		args = append(args, areaMax)
		countArgs = append(countArgs, areaMax)
	}

	// 添加排序
	if sortField, ok := filters["sortField"].(string); ok && sortField != "" {
		var dbField string
		switch sortField {
		case "id":
			dbField = "id"
		case "area":
			dbField = "house_area"
		case "fee":
			dbField = "fare_sum"
		default:
			dbField = "id"
		}

		sortOrder := "ASC"
		if order, ok := filters["sortOrder"].(string); ok && order == "desc" {
			sortOrder = "DESC"
		}

		query += fmt.Sprintf(" ORDER BY %s %s", dbField, sortOrder)
	} else {
		query += " ORDER BY id ASC"
	}

	// 添加分页
	query += " LIMIT ? OFFSET ?"
	args = append(args, pageSize, (page-1)*pageSize)

	// 打印最终的SQL查询
	fmt.Println("最终SQL查询:", query)
	fmt.Println("查询参数:", args)
	fmt.Println("计数SQL查询:", countQuery)
	fmt.Println("计数参数:", countArgs)

	// 执行查询获取总数
	var total int
	err := database.DB.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		fmt.Println("获取总数出错:", err)
		return nil, 0, err
	}
	fmt.Println("查询到的总数:", total)

	// 执行查询获取数据
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		fmt.Println("查询数据出错:", err)
		return nil, 0, err
	}
	defer rows.Close()

	// 解析结果
	var residents []Resident
	for rows.Next() {
		var r Resident
		err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.BlockNumber,
			&r.UnitNumber,
			&r.HouseNumber,
			&r.HouseArea,
			&r.FareSum,
			&r.Status,
			&r.CreatedAt,
			&r.UpdatedAt,
		)
		if err != nil {
			fmt.Println("解析行数据出错:", err)
			return nil, 0, err
		}

		// 设置状态文本
		if r.Status == 0 {
			r.StatusText = "欠费"
		} else {
			r.StatusText = "正常"
		}

		residents = append(residents, r)
	}

	fmt.Println("查询到的住户数量:", len(residents))
	if len(residents) > 0 {
		fmt.Println("第一个住户:", residents[0])
	}

	return residents, total, nil
}

// GetResidentByID 通过ID获取住户
func GetResidentByID(id int) (*Resident, error) {
	query := `SELECT id, name, block_number, unit_number, house_number, house_area, fare_sum, status, created_at, updated_at FROM residents WHERE id = ?`

	var r Resident
	err := database.DB.QueryRow(query, id).Scan(
		&r.ID,
		&r.Name,
		&r.BlockNumber,
		&r.UnitNumber,
		&r.HouseNumber,
		&r.HouseArea,
		&r.FareSum,
		&r.Status,
		&r.CreatedAt,
		&r.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 住户不存在
		}
		return nil, err
	}

	// 设置状态文本
	if r.Status == 0 {
		r.StatusText = "欠费"
	} else {
		r.StatusText = "正常"
	}

	return &r, nil
}

// CreateResident 创建住户
func CreateResident(req ResidentRequest) (int, error) {
	fmt.Println("接收到创建住户请求:", req)

	// 确保楼栋和单元格式正确
	blockNumber := req.BlockNumber
	if !strings.HasSuffix(blockNumber, "栋") {
		blockNumber = blockNumber + "栋"
	}

	unitNumber := req.UnitNumber
	if !strings.HasSuffix(unitNumber, "单元") {
		unitNumber = unitNumber + "单元"
	}

	query := `INSERT INTO residents (name, block_number, unit_number, house_number, house_area, fare_sum, status) VALUES (?, ?, ?, ?, ?, ?, ?)`
	fmt.Println("SQL查询:", query)
	fmt.Println("参数:", req.Name, blockNumber, unitNumber, req.HouseNumber, req.HouseArea, req.FareSum, req.Status)

	result, err := database.DB.Exec(
		query,
		req.Name,
		blockNumber,
		unitNumber,
		req.HouseNumber,
		req.HouseArea,
		req.FareSum,
		req.Status,
	)
	if err != nil {
		fmt.Println("创建住户失败:", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("获取新创建的住户ID失败:", err)
		return 0, err
	}

	fmt.Println("住户创建成功，ID:", id)
	return int(id), nil
}

// UpdateResident 更新住户
func UpdateResident(id int, req ResidentRequest) error {
	fmt.Println("接收到更新住户请求, ID:", id, "数据:", req)

	// 确保楼栋和单元格式正确
	blockNumber := req.BlockNumber
	if !strings.HasSuffix(blockNumber, "栋") {
		blockNumber = blockNumber + "栋"
	}

	unitNumber := req.UnitNumber
	if !strings.HasSuffix(unitNumber, "单元") {
		unitNumber = unitNumber + "单元"
	}

	query := `UPDATE residents SET name = ?, block_number = ?, unit_number = ?, house_number = ?, house_area = ?, fare_sum = ?, status = ? WHERE id = ?`
	fmt.Println("SQL查询:", query)
	fmt.Println("参数:", req.Name, blockNumber, unitNumber, req.HouseNumber, req.HouseArea, req.FareSum, req.Status, id)

	_, err := database.DB.Exec(
		query,
		req.Name,
		blockNumber,
		unitNumber,
		req.HouseNumber,
		req.HouseArea,
		req.FareSum,
		req.Status,
		id,
	)

	if err != nil {
		fmt.Println("更新住户失败:", err)
		return err
	}

	fmt.Println("住户更新成功, ID:", id)
	return nil
}

// DeleteResident 删除住户
func DeleteResident(id int) error {
	query := `DELETE FROM residents WHERE id = ?`

	_, err := database.DB.Exec(query, id)

	return err
}
