# Art Design Pro 权限系统对接指南（Golang 后端实现）

## 一、权限系统架构概览

Art Design Pro 支持**两种权限模式**：

### 1. 前端模式（Frontend Mode）
- 权限由前端路由配置控制
- 适合演示、小型项目
- 配置：`VITE_ACCESS_MODE = frontend`

### 2. 后端模式（Backend Mode）⭐️**推荐用于生产**
- 权限由后端 API 返回的菜单数据控制
- 适合企业级应用
- 配置：`VITE_ACCESS_MODE = backend`

---

## 二、权限控制的三个层级

### 1. **路由级权限**（基于角色 roles）
控制用户能访问哪些页面/菜单

### 2. **操作级权限**（基于权限标识 authMark）
控制页面内的按钮操作（增删改查）

### 3. **数据级权限**（后端实现）
控制用户能看到哪些数据（如：只能看自己部门的数据）

---

## 三、前端需要的 API 接口

### 1. 登录接口

**请求**
```http
POST /api/auth/login
Content-Type: application/json

{
  "userName": "admin",
  "password": "123456"
}
```

**响应**
```json
{
  "code": 200,
  "msg": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "refresh_token_string"
  }
}
```

### 2. 获取用户信息接口

**请求**
```http
GET /api/user/info
Authorization: Bearer {token}
```

**响应**
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "userId": 1,
    "userName": "admin",
    "email": "admin@example.com",
    "avatar": "https://example.com/avatar.jpg",
    "roles": ["R_SUPER", "R_ADMIN"],
    "buttons": ["add", "edit", "delete", "export"]
  }
}
```

**字段说明**：
- `roles`: 用户角色列表，用于**路由级权限**控制
- `buttons`: 用户拥有的操作权限标识，用于**操作级权限**控制

### 3. 获取菜单列表接口（后端模式核心）

**请求**
```http
GET /api/v3/system/menus
Authorization: Bearer {token}
```

**响应**（返回用户有权限的菜单树）
```json
{
  "code": 200,
  "msg": "success",
  "data": [
    {
      "id": 1,
      "path": "/dashboard",
      "name": "Dashboard",
      "component": "/dashboard/index",
      "meta": {
        "title": "工作台",
        "icon": "ri:dashboard-line",
        "keepAlive": true,
        "roles": ["R_SUPER", "R_ADMIN", "R_USER"]
      }
    },
    {
      "id": 2,
      "path": "/system",
      "name": "System",
      "component": "/index/index",
      "meta": {
        "title": "系统管理",
        "icon": "ri:settings-line",
        "roles": ["R_SUPER", "R_ADMIN"]
      },
      "children": [
        {
          "id": 21,
          "path": "user",
          "name": "User",
          "component": "/system/user",
          "meta": {
            "title": "用户管理",
            "icon": "ri:user-line",
            "keepAlive": true,
            "roles": ["R_SUPER", "R_ADMIN"],
            "authList": [
              { "title": "新增", "authMark": "add" },
              { "title": "编辑", "authMark": "edit" },
              { "title": "删除", "authMark": "delete" },
              { "title": "导出", "authMark": "export" }
            ]
          }
        },
        {
          "id": 22,
          "path": "role",
          "name": "Role",
          "component": "/system/role",
          "meta": {
            "title": "角色管理",
            "icon": "ri:user-settings-line",
            "keepAlive": true,
            "roles": ["R_SUPER"],
            "authList": [
              { "title": "新增", "authMark": "add" },
              { "title": "编辑", "authMark": "edit" },
              { "title": "删除", "authMark": "delete" }
            ]
          }
        }
      ]
    }
  ]
}
```

---

## 四、菜单数据结构详解

### 菜单项字段说明

```typescript
{
  "id": 1,                          // 菜单ID（可选，建议添加）
  "path": "/system/user",           // 路由路径
  "name": "User",                   // 路由名称（唯一标识）
  "component": "/system/user",      // 组件路径（相对于 src/views）
  "redirect": "/system/user/list",  // 重定向路径（可选）
  "meta": {
    "title": "用户管理",            // 菜单标题（支持i18n key）
    "icon": "ri:user-line",         // 图标（Iconify格式）
    "keepAlive": true,              // 是否缓存页面
    "isHide": false,                // 是否在菜单中隐藏
    "isHideTab": false,             // 是否在标签页中隐藏
    "fixedTab": false,              // 是否固定标签页
    "roles": ["R_SUPER", "R_ADMIN"], // 角色权限（路由级）
    "authList": [                   // 操作权限（按钮级）
      { "title": "新增", "authMark": "add" },
      { "title": "编辑", "authMark": "edit" },
      { "title": "删除", "authMark": "delete" }
    ],
    "link": "https://example.com",  // 外部链接（可选）
    "isIframe": false,              // 是否为iframe页面
    "isFullPage": false             // 是否为全屏页面
  },
  "children": []                    // 子菜单（递归结构）
}
```

### 特殊路由说明

**1. 目录型菜单**（有子菜单，自己不可点击）
```json
{
  "path": "/system",
  "name": "System",
  "component": "/index/index",  // 或空字符串 ""
  "meta": { "title": "系统管理" },
  "children": [...]
}
```

**2. 外部链接**
```json
{
  "path": "/external",
  "name": "External",
  "meta": {
    "title": "外部链接",
    "link": "https://github.com"
  }
}
```

**3. Iframe 页面**
```json
{
  "path": "/iframe-page",
  "name": "IframePage",
  "meta": {
    "title": "内嵌页面",
    "isIframe": true,
    "link": "https://example.com/page"
  }
}
```

---

## 五、前端权限验证逻辑

### 1. 路由守卫验证流程

```
用户访问页面
    ↓
检查是否登录
    ↓
获取用户信息（roles）
    ↓
获取菜单列表（后端返回）
    ↓
根据 roles 过滤菜单
    ↓
动态注册路由
    ↓
验证当前路径权限
    ↓
允许访问 / 跳转首页 / 跳转404
```

### 2. 指令级权限控制

**v-roles 指令**（页面级）
```vue
<!-- 只有 R_SUPER 或 R_ADMIN 能看到 -->
<el-button v-roles="['R_SUPER', 'R_ADMIN']">管理员功能</el-button>
```

**v-auth 指令**（操作级）
```vue
<!-- 只有拥有 add 权限的用户能看到 -->
<el-button v-auth="'add'">新增</el-button>
<el-button v-auth="'edit'">编辑</el-button>
<el-button v-auth="'delete'">删除</el-button>
```

---

## 六、Golang 后端实现示例

### 项目结构

```
golang-admin-best/
├── main.go
├── config/
│   └── config.go
├── models/
│   ├── user.go
│   ├── role.go
│   ├── menu.go
│   └── permission.go
├── controllers/
│   ├── auth_controller.go
│   └── system_controller.go
├── middleware/
│   ├── auth.go
│   └── cors.go
├── services/
│   ├── auth_service.go
│   ├── user_service.go
│   └── menu_service.go
├── routes/
│   └── routes.go
└── utils/
    ├── jwt.go
    └── response.go
```

### 1. 数据库表设计

```sql
-- 用户表
CREATE TABLE `users` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `user_name` VARCHAR(50) UNIQUE NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `email` VARCHAR(100),
  `avatar` VARCHAR(255),
  `status` TINYINT DEFAULT 1 COMMENT '1:启用 0:禁用',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 角色表
CREATE TABLE `roles` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `role_code` VARCHAR(50) UNIQUE NOT NULL COMMENT '角色代码，如 R_SUPER',
  `role_name` VARCHAR(50) NOT NULL COMMENT '角色名称',
  `description` VARCHAR(255),
  `status` TINYINT DEFAULT 1,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 用户角色关联表
CREATE TABLE `user_roles` (
  `user_id` INT NOT NULL,
  `role_id` INT NOT NULL,
  PRIMARY KEY (`user_id`, `role_id`),
  FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`role_id`) REFERENCES `roles`(`id`) ON DELETE CASCADE
);

-- 菜单表
CREATE TABLE `menus` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `parent_id` INT DEFAULT 0 COMMENT '父菜单ID，0表示一级菜单',
  `path` VARCHAR(255) NOT NULL COMMENT '路由路径',
  `name` VARCHAR(100) NOT NULL COMMENT '路由名称',
  `component` VARCHAR(255) COMMENT '组件路径',
  `redirect` VARCHAR(255) COMMENT '重定向路径',
  `title` VARCHAR(100) NOT NULL COMMENT '菜单标题',
  `icon` VARCHAR(100) COMMENT '图标',
  `keep_alive` TINYINT DEFAULT 0 COMMENT '是否缓存',
  `is_hide` TINYINT DEFAULT 0 COMMENT '是否隐藏',
  `is_iframe` TINYINT DEFAULT 0 COMMENT '是否iframe',
  `link` VARCHAR(255) COMMENT '外部链接',
  `sort` INT DEFAULT 0 COMMENT '排序',
  `status` TINYINT DEFAULT 1,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 角色菜单关联表
CREATE TABLE `role_menus` (
  `role_id` INT NOT NULL,
  `menu_id` INT NOT NULL,
  PRIMARY KEY (`role_id`, `menu_id`),
  FOREIGN KEY (`role_id`) REFERENCES `roles`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`menu_id`) REFERENCES `menus`(`id`) ON DELETE CASCADE
);

-- 权限表（按钮级权限）
CREATE TABLE `permissions` (
  `id` INT PRIMARY KEY AUTO_INCREMENT,
  `menu_id` INT NOT NULL COMMENT '所属菜单',
  `title` VARCHAR(50) NOT NULL COMMENT '权限名称',
  `auth_mark` VARCHAR(50) NOT NULL COMMENT '权限标识，如 add, edit, delete',
  FOREIGN KEY (`menu_id`) REFERENCES `menus`(`id`) ON DELETE CASCADE
);

-- 角色权限关联表
CREATE TABLE `role_permissions` (
  `role_id` INT NOT NULL,
  `permission_id` INT NOT NULL,
  PRIMARY KEY (`role_id`, `permission_id`),
  FOREIGN KEY (`role_id`) REFERENCES `roles`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`permission_id`) REFERENCES `permissions`(`id`) ON DELETE CASCADE
);
```

### 2. 初始化数据

```sql
-- 插入角色
INSERT INTO `roles` (`role_code`, `role_name`, `description`) VALUES
('R_SUPER', '超级管理员', '拥有所有权限'),
('R_ADMIN', '管理员', '拥有大部分权限'),
('R_USER', '普通用户', '基础权限');

-- 插入测试用户（密码: 123456，需要加密）
INSERT INTO `users` (`user_name`, `password`, `email`, `avatar`) VALUES
('admin', '$2a$10$...', 'admin@example.com', 'https://example.com/avatar.jpg');

-- 关联用户和角色
INSERT INTO `user_roles` (`user_id`, `role_id`) VALUES (1, 1);

-- 插入菜单（示例）
INSERT INTO `menus` (`parent_id`, `path`, `name`, `component`, `title`, `icon`, `sort`) VALUES
(0, '/dashboard', 'Dashboard', '/dashboard/index', '工作台', 'ri:dashboard-line', 1),
(0, '/system', 'System', '/index/index', '系统管理', 'ri:settings-line', 2);

INSERT INTO `menus` (`parent_id`, `path`, `name`, `component`, `title`, `icon`, `sort`) VALUES
(2, 'user', 'User', '/system/user', '用户管理', 'ri:user-line', 1),
(2, 'role', 'Role', '/system/role', '角色管理', 'ri:user-settings-line', 2);

-- 关联角色和菜单
INSERT INTO `role_menus` (`role_id`, `menu_id`) VALUES
(1, 1), (1, 2), (1, 3), (1, 4),  -- 超级管理员拥有所有菜单
(2, 1), (2, 2), (2, 3),          -- 管理员部分菜单
(3, 1);                          -- 普通用户只有工作台

-- 插入权限
INSERT INTO `permissions` (`menu_id`, `title`, `auth_mark`) VALUES
(3, '新增', 'add'),
(3, '编辑', 'edit'),
(3, '删除', 'delete'),
(3, '导出', 'export');

-- 关联角色和权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`) VALUES
(1, 1), (1, 2), (1, 3), (1, 4),  -- 超级管理员所有权限
(2, 1), (2, 2);                  -- 管理员部分权限
```

### 3. 核心代码实现

详见下一个文件...
