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

// allRoles 所有角色代码（菜单未显式限制角色时默认对全部角色可见）
var allRoles = []string{"R_SUPER", "R_ADMIN", "R_USER"}

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
		roleIDByCode := make(map[string]uint, len(roles))
		for _, r := range roles {
			roleIDByCode[r.RoleCode] = r.ID
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

		// 4. 递归插入菜单树（数据见 menu_seed.go），并按角色关联菜单/权限
		if err := insertMenus(tx, menuSeeds(), 0, allRoles, roleIDByCode); err != nil {
			return err
		}
		log.Println("   → Menus & permissions created")

		return nil
	})
}

// insertMenus 递归插入菜单，处理 parent_id、排序、按钮权限，并按“有效角色”关联 role_menus/role_permissions。
// inherited 为父级下放的可见角色；子菜单的有效角色 = inherited ∩ 自身 Roles（自身未设则继承 inherited），避免出现子菜单可见而父菜单不可见的孤儿。
func insertMenus(tx *gorm.DB, items []seedMenu, parentID uint, inherited []string, roleIDByCode map[string]uint) error {
	for i, it := range items {
		m := &model.Menu{
			ParentID:      parentID,
			Path:          it.Path,
			Name:          it.Name,
			Component:     it.Component,
			Redirect:      it.Redirect,
			Title:         it.Title,
			Icon:          it.Icon,
			KeepAlive:     b2i(it.KeepAlive),
			IsHide:        b2i(it.IsHide),
			IsHideTab:     b2i(it.IsHideTab),
			IsIframe:      b2i(it.IsIframe),
			IsFullPage:    b2i(it.IsFullPage),
			FixedTab:      b2i(it.FixedTab),
			Link:          it.Link,
			ActivePath:    it.ActivePath,
			ShowTextBadge: it.ShowTextBadge,
			Sort:          i + 1,
			Status:        1,
		}
		if err := tx.Create(m).Error; err != nil {
			return err
		}

		// 计算有效角色
		eff := inherited
		if len(it.Roles) > 0 {
			eff = intersect(inherited, it.Roles)
		}

		// 关联 role_menus
		for _, rc := range eff {
			if rid, ok := roleIDByCode[rc]; ok {
				if err := tx.Exec("INSERT INTO role_menus (role_id, menu_id) VALUES (?, ?)", rid, m.ID).Error; err != nil {
					return err
				}
			}
		}

		// 按钮权限 + role_permissions
		for _, a := range it.Auth {
			p := &model.Permission{MenuID: m.ID, Title: a.Title, AuthMark: a.AuthMark}
			if err := tx.Create(p).Error; err != nil {
				return err
			}
			for _, rc := range eff {
				if rid, ok := roleIDByCode[rc]; ok {
					if err := tx.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", rid, p.ID).Error; err != nil {
						return err
					}
				}
			}
		}

		// 递归子菜单
		if len(it.Children) > 0 {
			if err := insertMenus(tx, it.Children, m.ID, eff, roleIDByCode); err != nil {
				return err
			}
		}
	}
	return nil
}

// b2i 布尔转 int8（1/0）
func b2i(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

// intersect 求交集，保留 a 的顺序
func intersect(a, b []string) []string {
	set := make(map[string]struct{}, len(b))
	for _, x := range b {
		set[x] = struct{}{}
	}
	out := make([]string, 0, len(a))
	for _, x := range a {
		if _, ok := set[x]; ok {
			out = append(out, x)
		}
	}
	return out
}
