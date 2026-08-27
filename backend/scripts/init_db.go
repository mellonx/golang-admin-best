package main

import (
	"golang-admin-best/internal/config"
	"golang-admin-best/internal/model"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// 加载配置
	if err := config.Load(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	dbConfig := config.AppConfig.DB["default"]
	if dbConfig == nil {
		log.Fatal("Default database config not found")
	}

	// 1. 连接 MySQL（不指定数据库）
	db, err := gorm.Open(mysql.Open(dbConfig.GetDSNWithoutDB()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}

	// 2. 创建数据库
	dbName := dbConfig.DBName
	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)
	if err := db.Exec(createDBSQL).Error; err != nil {
		log.Fatal("Failed to create database:", err)
	}
	log.Printf("✅ Database '%s' created or already exists\n", dbName)

	// 3. 切换到目标数据库
	db, err = gorm.Open(mysql.Open(dbConfig.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 4. 自动迁移表结构
	log.Println("🔧 Migrating tables...")
	if err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Menu{},
		&model.Permission{},
	); err != nil {
		log.Fatal("Failed to migrate tables:", err)
	}
	log.Println("✅ Tables migrated successfully")

	// 5. 初始化数据
	log.Println("📦 Initializing data...")
	if err := initData(db); err != nil {
		log.Fatal("Failed to initialize data:", err)
	}
	log.Println("✅ Data initialized successfully")

	log.Println("🎉 Database setup complete!")
}

func initData(db *gorm.DB) error {
	// 检查是否已初始化
	var count int64
	db.Model(&model.Role{}).Count(&count)
	if count > 0 {
		log.Println("⚠️  Data already exists, skipping initialization")
		return nil
	}

	// 开启事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. 插入角色
		roles := []*model.Role{
			{RoleCode: "R_SUPER", RoleName: "超级管理员", Description: "拥有所有权限"},
			{RoleCode: "R_ADMIN", RoleName: "管理员", Description: "拥有大部分权限"},
			{RoleCode: "R_USER", RoleName: "普通用户", Description: "基础权限"},
		}
		if err := tx.Create(&roles).Error; err != nil {
			return err
		}
		log.Println("   → Roles created")

		// 2. 插入用户（密码均为 123456，用户名对应前端登录页的快捷账号）
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		users := []*model.User{
			{UserName: "Super", Password: string(hashedPassword), Email: "super@example.com", Avatar: "https://avatars.githubusercontent.com/u/1?v=4", Status: 1},
			{UserName: "Admin", Password: string(hashedPassword), Email: "admin@example.com", Avatar: "https://avatars.githubusercontent.com/u/2?v=4", Status: 1},
			{UserName: "User", Password: string(hashedPassword), Email: "user@example.com", Avatar: "https://avatars.githubusercontent.com/u/3?v=4", Status: 1},
		}
		if err := tx.Create(&users).Error; err != nil {
			return err
		}
		log.Println("   → Users created: Super/Admin/User (password: 123456)")

		// 3. 关联用户和角色（Super→R_SUPER, Admin→R_ADMIN, User→R_USER）
		userRoles := []struct{ userID, roleID uint }{
			{users[0].ID, roles[0].ID},
			{users[1].ID, roles[1].ID},
			{users[2].ID, roles[2].ID},
		}
		for _, ur := range userRoles {
			if err := tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", ur.userID, ur.roleID).Error; err != nil {
				return err
			}
		}

		// 4. 插入顶级目录菜单（component 使用 /index/index 布局，title 用 i18n key）
		dashboard := &model.Menu{ParentID: 0, Path: "/dashboard", Name: "Dashboard", Component: "/index/index", Title: "menus.dashboard.title", Icon: "ri:pie-chart-line", Sort: 1}
		system := &model.Menu{ParentID: 0, Path: "/system", Name: "System", Component: "/index/index", Title: "menus.system.title", Icon: "ri:user-3-line", Sort: 2}
		if err := tx.Create([]*model.Menu{dashboard, system}).Error; err != nil {
			return err
		}

		// 5. 插入子菜单（component 对应 src/views 下真实组件）
		console := &model.Menu{ParentID: dashboard.ID, Path: "console", Name: "Console", Component: "/dashboard/console", Title: "menus.dashboard.console", Icon: "ri:home-smile-2-line", Sort: 1}
		analysis := &model.Menu{ParentID: dashboard.ID, Path: "analysis", Name: "Analysis", Component: "/dashboard/analysis", Title: "menus.dashboard.analysis", Icon: "ri:align-item-bottom-line", Sort: 2}
		ecommerce := &model.Menu{ParentID: dashboard.ID, Path: "ecommerce", Name: "Ecommerce", Component: "/dashboard/ecommerce", Title: "menus.dashboard.ecommerce", Icon: "ri:bar-chart-box-line", Sort: 3}
		userMenu := &model.Menu{ParentID: system.ID, Path: "user", Name: "User", Component: "/system/user", Title: "menus.system.user", Icon: "ri:user-line", KeepAlive: 1, Sort: 1}
		roleMenu := &model.Menu{ParentID: system.ID, Path: "role", Name: "Role", Component: "/system/role", Title: "menus.system.role", Icon: "ri:user-settings-line", KeepAlive: 1, Sort: 2}
		if err := tx.Create([]*model.Menu{console, analysis, ecommerce, userMenu, roleMenu}).Error; err != nil {
			return err
		}
		log.Println("   → Menus created")

		// 6. 关联角色和菜单
		//    R_SUPER: 全部；R_ADMIN: 除角色管理外；R_USER: 仅工作台+控制台
		roleMenus := []struct{ roleID, menuID uint }{
			// 超级管理员：所有菜单
			{roles[0].ID, dashboard.ID}, {roles[0].ID, console.ID}, {roles[0].ID, analysis.ID}, {roles[0].ID, ecommerce.ID},
			{roles[0].ID, system.ID}, {roles[0].ID, userMenu.ID}, {roles[0].ID, roleMenu.ID},
			// 管理员：工作台全部 + 系统管理(用户，不含角色)
			{roles[1].ID, dashboard.ID}, {roles[1].ID, console.ID}, {roles[1].ID, analysis.ID}, {roles[1].ID, ecommerce.ID},
			{roles[1].ID, system.ID}, {roles[1].ID, userMenu.ID},
			// 普通用户：仅工作台 + 控制台
			{roles[2].ID, dashboard.ID}, {roles[2].ID, console.ID},
		}
		for _, rm := range roleMenus {
			tx.Exec("INSERT INTO role_menus (role_id, menu_id) VALUES (?, ?)", rm.roleID, rm.menuID)
		}

		// 7. 插入按钮权限（挂在用户管理菜单下）
		perms := []*model.Permission{
			{MenuID: userMenu.ID, Title: "新增", AuthMark: "add"},
			{MenuID: userMenu.ID, Title: "编辑", AuthMark: "edit"},
			{MenuID: userMenu.ID, Title: "删除", AuthMark: "delete"},
			{MenuID: userMenu.ID, Title: "导出", AuthMark: "export"},
		}
		if err := tx.Create(&perms).Error; err != nil {
			return err
		}
		log.Println("   → Permissions created")

		// 8. 关联角色和权限
		rolePerms := []struct{ roleID, permID uint }{
			{roles[0].ID, perms[0].ID}, {roles[0].ID, perms[1].ID}, {roles[0].ID, perms[2].ID}, {roles[0].ID, perms[3].ID},
			{roles[1].ID, perms[0].ID}, {roles[1].ID, perms[1].ID},
		}
		for _, rp := range rolePerms {
			tx.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", rp.roleID, rp.permID)
		}

		return nil
	})
}
