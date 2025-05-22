# Community Backward

一个基于Go的Web应用程序，使用Gin框架和MySQL数据库。

## 功能

- 用户认证（登录）
- JWT令牌认证
- 仪表盘数据展示

## 技术栈

- Go 1.x
- Gin Web框架
- MySQL数据库
- JWT认证

## 项目结构

```
community-backward/
├── config/             # 应用配置
├── database/           # 数据库连接和模式
├── handlers/           # HTTP处理函数
├── middleware/         # 中间件
├── models/             # 数据模型
├── routes/             # 路由配置
├── utils/              # 工具函数
├── .env                # 环境变量
├── go.mod              # Go模块文件
├── go.sum              # 依赖校验文件
└── main.go             # 应用入口点
```

## 数据库设置

1. 确保MySQL服务器正在运行
2. 创建数据库和表：

```sql
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

-- 插入单个用户数据（密码为123456的MD5哈希值）
-- 123456的MD5哈希值为e10adc3949ba59abbe56e057f20f883e
INSERT INTO users (username, password, email) VALUES
('admin', 'e10adc3949ba59abbe56e057f20f883e', 'admin@example.com')
ON DUPLICATE KEY UPDATE password = VALUES(password), email = VALUES(email);
```

或者，您可以运行`database/schema.sql`文件：

```bash
mysql -u root -p < database/schema.sql
```

## 环境变量

在`.env`文件中配置以下环境变量：

```
# 服务器配置
PORT=3000
GIN_MODE=debug

# JWT配置
JWT_SECRET=your_jwt_secret_key

# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=community
```

## 运行应用

1. 安装依赖：

```bash
go mod download
```

2. 构建应用：

```bash
go build
```

3. 运行应用：

```bash
./community-backward
```

## API端点

### 认证

- `POST /api/auth/login` - 用户登录
  - 请求体: `{ "username": "admin", "password": "123456", "remember": true }`
  - 响应: `{ "success": true, "token": "jwt_token", "message": "登录成功" }`



### 仪表盘

- `GET /api/dashboard` - 获取仪表盘数据（需要认证）
  - 请求头: `Authorization: Bearer jwt_token`
  - 响应: `{ "success": true, "data": { ... } }`
