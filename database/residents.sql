-- 创建住户表
CREATE TABLE IF NOT EXISTS residents (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    block_number VARCHAR(10) NOT NULL,
    unit_number VARCHAR(10) NOT NULL,
    house_number VARCHAR(10) NOT NULL,
    house_area DECIMAL(10,2) NOT NULL,
    fare_sum DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    status TINYINT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入示例数据
INSERT INTO residents (name, block_number, unit_number, house_number, house_area, fare_sum, status) VALUES
('张三', '1栋', '1单元', '101', 85.50, 1200.00, 1),
('李四', '1栋', '1单元', '102', 120.30, 1800.00, 1),
('王五', '1栋', '2单元', '201', 95.00, 0.00, 0),
('赵六', '2栋', '1单元', '101', 75.80, 1000.00, 1),
('钱七', '2栋', '2单元', '202', 110.50, 0.00, 0),
('孙八', '3栋', '1单元', '101', 90.20, 1350.00, 1);
