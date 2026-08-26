# Art Design Pro - Backend API

Golang 后端服务，为 Art Design Pro 前端提供权限管理和业务接口。

## 技术栈

- Golang 1.21+
- Gin Web Framework
- GORM (ORM)
- MySQL
- JWT 认证

## 项目结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go           # 入口程序
├── internal/                 # 私有代码
│   ├── handler/              # HTTP 处理器
│   ├── service/              # 业务逻辑
│   ├── model/                # 数据模型
│   ├── middleware/           # 中间件
│   └── database/             # 数据库连接
├── pkg/                      # 公共工具
│   ├── response/             # 统一响应格式
│   └── utils/                # 工具函数
├── configs/                  # 配置文件
├── scripts/                  # SQL 脚本
└── go.mod
```

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置数据库

在 MySQL 中创建数据库：

```sql
CREATE DATABASE art_design_pro CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

执行初始化脚本（待创建）：

```bash
mysql -u root -p art_design_pro < scripts/init.sql
```

### 3. 配置环境变量

复制配置文件并修改：

```bash
cp configs/config.example.yaml configs/config.yaml
# 编辑 config.yaml 填入数据库连接信息
```

### 4. 运行

```bash
go run cmd/server/main.go
```

API 将运行在 http://localhost:8080

## 核心接口

详见 [../docs/backend-guide.md](../docs/backend-guide.md)

### 登录接口

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "123456"
}
```

### 获取用户信息

```http
GET /api/user/info
Authorization: Bearer {token}
```

### 获取菜单树

```http
GET /api/v3/system/menus
Authorization: Bearer {token}
```

## 开发说明

- 遵循 Go 标准项目布局
- `internal/` 包不对外暴露
- `pkg/` 可被其他项目引用
- 使用 GORM 的 AutoMigrate 自动创建表结构

## License

MIT License
