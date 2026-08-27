package main

import (
	"art-design-pro-api/internal/config"
	"art-design-pro-api/internal/model"
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

	// 1. 连接 MySQL（不指定数据库）
	db, err := gorm.Open(mysql.Open(config.AppConfig.DB.GetDSNWithoutDB()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}

	// 2. 创建数据库
	dbName := config.AppConfig.DB.DBName
	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)
	if err := db.Exec(createDBSQL).Error; err != nil {
		log.Fatal("Failed to create database:", err)
	}
	log.Printf("✅ Database '%s' created or already exists\n", dbName)

	// 3. 切换到目标数据库
	db, err = gorm.Open(mysql.Open(config.AppConfig.DB.GetDSN()), &gorm.Config{
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

		// 2. 插入用户（密码: 123456）
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		user := &model.User{
			UserName: "admin",
			Password: string(hashedPassword),
			Email:    "admin@example.com",
			Avatar:   "https://avatars.githubusercontent.com/u/1?v=4",
			Status:   1,
		}
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		log.Println("   → User 'admin' created (password: 123456)")

		// 3. 关联用户和角色（admin 拥有超级管理员角色）
		if err := tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", user.ID, roles[0].ID).Error; err != nil {
			return err
		}

		// 4. 插入菜单
		menus := []*model.Menu{
			{ParentID: 0, Path: "/dashboard", Name: "Dashboard", Component: "/dashboard/index", Title: "工作台", Icon: "ri:dashboard-line", Sort: 1},
			{ParentID: 0, Path: "/system", Name: "System", Component: "/index/index", Title: "系统管理", Icon: "ri:settings-line", Sort: 2},
		}
		if err := tx.Create(&menus).Error; err != nil {
			return err
		}

		subMenus := []*model.Menu{
			{ParentID: menus[1].ID, Path: "user", Name: "User", Component: "/system/user", Title: "用户管理", Icon: "ri:user-line", Sort: 1},
			{ParentID: menus[1].ID, Path: "role", Name: "Role", Component: "/system/role", Title: "角色管理", Icon: "ri:user-settings-line", Sort: 2},
		}
		if err := tx.Create(&subMenus).Error; err != nil {
			return err
		}
		log.Println("   → Menus created")

		// 5. 关联角色和菜单
		roleMenus := []struct{ roleID, menuID uint }{
			{roles[0].ID, menus[0].ID}, {roles[0].ID, menus[1].ID}, {roles[0].ID, subMenus[0].ID}, {roles[0].ID, subMenus[1].ID},
			{roles[1].ID, menus[0].ID}, {roles[1].ID, menus[1].ID}, {roles[1].ID, subMenus[0].ID},
			{roles[2].ID, menus[0].ID},
		}
		for _, rm := range roleMenus {
			tx.Exec("INSERT INTO role_menus (role_id, menu_id) VALUES (?, ?)", rm.roleID, rm.menuID)
		}

		// 6. 插入权限（用户管理的按钮权限）
		perms := []*model.Permission{
			{MenuID: subMenus[0].ID, Title: "新增", AuthMark: "add"},
			{MenuID: subMenus[0].ID, Title: "编辑", AuthMark: "edit"},
			{MenuID: subMenus[0].ID, Title: "删除", AuthMark: "delete"},
			{MenuID: subMenus[0].ID, Title: "导出", AuthMark: "export"},
		}
		if err := tx.Create(&perms).Error; err != nil {
			return err
		}
		log.Println("   → Permissions created")

		// 7. 关联角色和权限
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
