# Art Design Pro - Backend API

Golang 后端服务，为 Art Design Pro 前端提供权限管理和业务接口。

## 技术栈

- Golang 1.21+
- Gin Web Framework
- GORM (ORM)
- MySQL
- JWT 认证

## 项目结构

分层架构：`handler`（HTTP）→ `service`（业务逻辑）→ `repository`（数据访问）→ `model`（数据模型）

```
backend/
├── cmd/
│   └── server/
│       └── main.go           # 入口程序（依赖注入 + 路由注册）
├── internal/                 # 私有代码
│   ├── config/               # 配置加载（.env / 环境变量）
│   ├── database/             # 数据库连接（支持多命名连接）
│   ├── model/                # 数据模型 + DTO
│   ├── repository/           # 数据访问层（接口 + 实现）
│   ├── service/              # 业务逻辑层
│   ├── handler/              # HTTP 处理器
│   ├── router/               # 路由注册（Setup + Handlers 聚合）
│   └── middleware/           # 中间件（JWT 认证 / CORS）
├── pkg/                      # 公共工具
│   ├── response/             # 统一响应格式
│   └── utils/                # JWT 工具
├── scripts/
│   └── init_db.go            # 建库建表 + 初始化数据
├── .env.example              # 配置模板
└── go.mod
```

## 配置管理

配置优先从 `.env` 文件读取，不存在的项回退到系统环境变量。

```bash
cp .env.example .env
# 编辑 .env 填入数据库连接、JWT 密钥等
```

### 多数据库连接

`internal/config/config.go` 中的 `DB` 是一个 `map[string]*DBConfig`，默认包含 `default` 连接。
如需从库、多租户等场景，在该 map 中添加命名连接，通过 `database.Get("连接名")` 获取对应实例。

```go
repo := repository.NewRepository(database.Get())          // 默认连接
slaveRepo := repository.NewRepository(database.Get("slave")) // 从库连接
```

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置并初始化数据库

确保 `.env` 中数据库配置正确，然后运行初始化脚本（自动建库、建表、插入初始数据）：

```bash
go run scripts/init_db.go
```

这会创建 `golang_admin_best` 数据库、7 张表，以及三个初始账号（密码均为 `123456`），对应前端登录页的快捷登录：

| 用户名 | 角色 | 权限范围 |
|--------|------|---------|
| `Super` | R_SUPER | 全部菜单 + 全部按钮权限 |
| `Admin` | R_ADMIN | 工作台 + 用户管理（无角色管理），按钮 add/edit |
| `User` | R_USER | 仅工作台/控制台 |

### 3. 运行

```bash
go run cmd/server/main.go
```

API 将运行在 http://localhost:8080

## 核心接口

详见 [../docs/backend-guide.md](../docs/backend-guide.md)

> **认证格式**：受保护接口在请求头携带 `Authorization: {token}`（前端当前用法），
> 中间件同时兼容标准的 `Authorization: Bearer {token}` 格式。

### 登录接口（公开）

```http
POST /api/auth/login
Content-Type: application/json

{
  "userName": "Super",
  "password": "123456"
}
```

响应：`{ code, msg, data: { token, refreshToken } }`

### 获取用户信息（需认证）

```http
GET /api/user/info
Authorization: {token}
```

响应 `data`：`{ userId, userName, email, avatar, roles[], buttons[] }`
- `roles`：角色代码（如 `R_SUPER`），用于路由级权限
- `buttons`：按钮权限标识（如 `add`/`edit`），用于操作级权限

### 获取菜单树（需认证）

```http
GET /api/v3/system/menus
Authorization: {token}
```

响应 `data`：根据用户角色过滤的菜单树，含 `meta.authList` 按钮权限。

## 开发说明

- 遵循 Go 标准项目布局
- `internal/` 包不对外暴露，`pkg/` 可被其他项目引用
- 分层职责：handler 只处理 HTTP，业务逻辑在 service，数据访问在 repository
- Repository 层为接口抽象，便于替换实现和单元测试
- 响应统一返回 HTTP 200，业务状态码在 body 的 `code` 字段（对应前端 ApiStatus）

## 前端对接

修改前端 `frontend/.env`：

```bash
VITE_ACCESS_MODE = backend        # 切换到后端权限模式
```

确认前端 API 代理指向后端服务（默认 http://localhost:8080）。

## License

MIT License
