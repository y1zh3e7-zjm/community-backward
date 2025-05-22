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


-- 插入一些示例物业费记录
-- 注意：请确保residents表中存在ID为1-6的记录，否则会出现外键约束错误
-- 如果residents表中的ID不是1-6，请修改下面的resident_id值以匹配实际存在的ID

-- 首先确保residents表中有数据
INSERT IGNORE INTO residents (id, name, block_number, unit_number, house_number, house_area, fare_sum, status) VALUES
(1, '张三', '1栋', '1单元', '101', 85.50, 1200.00, 1),
(2, '李四', '1栋', '1单元', '102', 120.30, 1800.00, 1),
(3, '王五', '1栋', '2单元', '201', 95.00, 0.00, 0),
(4, '赵六', '2栋', '1单元', '101', 75.80, 1000.00, 1),
(5, '钱七', '2栋', '2单元', '202', 110.50, 0.00, 0),
(6, '孙八', '3栋', '1单元', '101', 90.20, 1350.00, 1);

-- 然后插入物业费记录
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
