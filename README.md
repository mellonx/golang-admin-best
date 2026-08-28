# Art Design Pro - Full Stack Solution

基于 MIT 协议的企业级管理后台全栈解决方案，前后端分离架构。

> 原项目：[Daymychen/art-design-pro](https://github.com/Daymychen/art-design-pro)  
> 本仓库为 monorepo 重构版本，包含前端（Vue 3）和后端（Golang）

## 📁 项目结构

```
art-design-pro/
├── frontend/          # Vue 3 前端项目
│   ├── src/          # 源代码
│   ├── public/       # 静态资源
│   ├── README.md     # 前端详细文档
│   └── package.json
├── backend/           # Golang 后端 API
│   ├── cmd/          # 入口程序
│   ├── internal/     # 内部包（config, database, model, repository, service, handler, router, middleware）
│   ├── pkg/          # 公共工具（response, utils）
│   ├── scripts/      # 建库建表 + 初始化数据
│   └── go.mod
└── docs/              # 项目文档
    ├── backend-guide.md    # 权限系统对接指南
    └── golang-code.md      # Golang 实现示例
```

## 🚀 快速开始

### 前端开发

```bash
cd frontend
pnpm install
pnpm dev
```

访问 http://localhost:3006

详细说明见 [frontend/README.md](frontend/README.md)

### 后端开发

```bash
cd backend
go mod download
cp .env.example .env        # 按需修改数据库连接、JWT 密钥
go run ./scripts            # 建库建表 + 初始化数据（含菜单/角色/账号）
go run cmd/server/main.go
```

API 运行在 http://localhost:8080

初始化后会创建数据库 `golang_admin_best` 与三个演示账号（密码均为 `123456`）：
`Super`（超级管理员）/ `Admin`（管理员）/ `User`（普通用户）。

详细说明见 [backend/README.md](backend/README.md) 与 [docs/backend-guide.md](docs/backend-guide.md)

## 🛠 技术栈

### 前端
- **框架**: Vue 3.5 + TypeScript
- **构建**: Vite 7
- **UI**: Tailwind CSS 4 + Element Plus 2.11
- **状态**: Pinia
- **路由**: Vue Router
- **权限**: 路由级 / 操作级 / 数据级三层权限控制

### 后端
- **语言**: Golang（module `golang-admin-best`）
- **框架**: Gin
- **ORM**: GORM（支持多命名数据库连接）
- **数据库**: MySQL（默认库名 `golang_admin_best`）
- **认证**: JWT
- **架构**: handler → service → repository → model 分层

## 📖 开发说明

### 前端权限模式

前端默认运行在**前端权限模式**（使用 Mock 数据），开箱即用。

### 后端权限模式（已实现）

前端已切换为后端模式（`frontend/.env` 中 `VITE_ACCESS_MODE = backend`），
菜单与权限由后端 API 提供。已实现接口：

| 接口 | 说明 |
|------|------|
| `POST /api/auth/login` | 登录，返回 token / refreshToken |
| `GET /api/user/info` | 用户信息（roles + buttons 权限） |
| `GET /api/v3/system/menus` | 按角色过滤的菜单树 |
| `GET /api/user/list` | 用户列表（分页） |
| `GET /api/role/list` | 角色列表（分页） |

菜单数据已与前端 11 个路由模块对齐，登录后按角色显示不同菜单。
详细 API 规范见 [backend/README.md](backend/README.md) 与 [docs/backend-guide.md](docs/backend-guide.md)。

## 📝 License

MIT License - 可自由用于商业项目

- ✅ 商业使用
- ✅ 修改源码
- ✅ 分发和销售
- ⚠️ 需保留原作者版权声明

详见 [LICENSE](LICENSE)

## 🙏 致谢

感谢 [SuperManTT](https://github.com/Daymychen) 创建的优秀开源项目 Art Design Pro。
