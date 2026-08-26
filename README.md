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
│   ├── internal/     # 内部包（handler, service, model, middleware）
│   ├── pkg/          # 公共工具
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
go run cmd/server/main.go
```

API 运行在 http://localhost:8080

详细对接方案见 [docs/backend-guide.md](docs/backend-guide.md)

## 🛠 技术栈

### 前端
- **框架**: Vue 3.5 + TypeScript
- **构建**: Vite 7
- **UI**: Tailwind CSS 4 + Element Plus 2.11
- **状态**: Pinia
- **路由**: Vue Router
- **权限**: 路由级 / 操作级 / 数据级三层权限控制

### 后端（待实现）
- **语言**: Golang 1.21+
- **框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **认证**: JWT

## 📖 开发说明

### 前端权限模式

前端默认运行在**前端权限模式**（使用 Mock 数据），开箱即用。

### 后端权限模式

1. 修改 `frontend/.env`：
   ```bash
   VITE_ACCESS_MODE = backend
   ```

2. 启动后端服务（需先实现，见 `docs/` 文档）

3. 后端需实现 3 个核心接口：
   - `POST /api/auth/login` - 登录
   - `GET /api/user/info` - 用户信息
   - `GET /api/v3/system/menus` - 菜单树

详细 API 规范见 [docs/backend-guide.md](docs/backend-guide.md)

## 📝 License

MIT License - 可自由用于商业项目

- ✅ 商业使用
- ✅ 修改源码
- ✅ 分发和销售
- ⚠️ 需保留原作者版权声明

详见 [LICENSE](LICENSE)

## 🙏 致谢

感谢 [SuperManTT](https://github.com/Daymychen) 创建的优秀开源项目 Art Design Pro。
