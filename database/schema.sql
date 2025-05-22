-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS community CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE community;

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(32) NOT NULL, -- MD5哈希值为32位
    email VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入单个用户数据（密码为hexiewuye123123的MD5哈希值）
-- 123456的MD5哈希值为e10adc3949ba59abbe56e057f20f883e
INSERT INTO users (username, password, email) VALUES
('admin', 'edafa832c7404ebd47234eddacad34f1', '2510733026@qq.com')
ON DUPLICATE KEY UPDATE password = VALUES(password), email = VALUES(email);

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
INSERT IGNORE INTO residents (id, name, block_number, unit_number, house_number, house_area, fare_sum, status) VALUES
(1, '张三', '1栋', '1单元', '101', 85.50, 1200.00, 1),
(2, '李四', '1栋', '1单元', '102', 120.30, 1800.00, 1),
(3, '王五', '1栋', '2单元', '201', 95.00, 0.00, 0),
(4, '赵六', '2栋', '1单元', '101', 75.80, 1000.00, 1),
(5, '钱七', '2栋', '2单元', '202', 110.50, 0.00, 0),
(6, '孙八', '3栋', '1单元', '101', 90.20, 1350.00, 1);

-- 创建物业费表
CREATE TABLE IF NOT EXISTS property_fees (
    id INT AUTO_INCREMENT PRIMARY KEY,
    resident_id INT NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    payment_date DATE NOT NULL,
    payment_method VARCHAR(20) NOT NULL,
    payment_status TINYINT NOT NULL DEFAULT 1,
    remark VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (resident_id) REFERENCES residents(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建物业费月度统计表
CREATE TABLE IF NOT EXISTS property_fee_monthly_stats (
    id INT AUTO_INCREMENT PRIMARY KEY,
    year INT NOT NULL,
    month INT NOT NULL,
    actual_income DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    planned_income DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY year_month_idx (year, month)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入2023年的示例数据
INSERT INTO property_fee_monthly_stats (year, month, actual_income, planned_income) VALUES
(2023, 1, 82000.00, 85000.00),
(2023, 2, 78500.00, 82000.00),
(2023, 3, 85300.00, 85000.00),
(2023, 4, 90100.00, 88000.00),
(2023, 5, 92500.00, 90000.00),
(2023, 6, 88700.00, 92000.00),
(2023, 7, 91200.00, 93000.00),
(2023, 8, 89500.00, 94000.00),
(2023, 9, 93800.00, 95000.00),
(2023, 10, 96200.00, 96000.00),
(2023, 11, 92300.00, 97000.00),
(2023, 12, 98500.00, 98000.00);

-- 插入2024年的示例数据
INSERT INTO property_fee_monthly_stats (year, month, actual_income, planned_income) VALUES
(2024, 1, 88000.00, 90000.00),
(2024, 2, 85500.00, 90000.00),
(2024, 3, 92300.00, 92000.00),
(2024, 4, 94100.00, 93000.00),
(2024, 5, 96500.00, 95000.00),
(2024, 6, 93700.00, 96000.00);

-- 插入一些示例物业费记录
-- 确保先有residents表中的数据，再插入property_fees表的数据
INSERT INTO property_fees (resident_id, amount, payment_date, payment_method, payment_status, remark) VALUES
(1, 1200.00, '2024-01-15', '微信支付', 1, '2024年1月物业费'),
(2, 1800.00, '2024-01-20', '支付宝', 1, '2024年1月物业费'),
(3, 950.00, '2024-02-05', '银行转账', 1, '2024年1-2月物业费'),
(4, 1000.00, '2024-02-10', '微信支付', 1, '2024年2月物业费'),
(5, 1105.00, '2024-03-15', '支付宝', 1, '2024年3月物业费'),
(6, 902.00, '2024-03-20', '现金', 1, '2024年3月物业费'),
(1, 1200.00, '2024-04-15', '微信支付', 1, '2024年4月物业费'),
(2, 1800.00, '2024-04-20', '支付宝', 1, '2024年4月物业费'),
(3, 950.00, '2024-05-05', '银行转账', 1, '2024年5月物业费'),
(4, 1000.00, '2024-05-10', '微信支付', 1, '2024年5月物业费'),
(5, 1105.00, '2024-06-15', '支付宝', 1, '2024年6月物业费'),
(6, 902.00, '2024-06-20', '现金', 1, '2024年6月物业费');
