package repository

import (
	"golang-admin-best/internal/model"

	"gorm.io/gorm"
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	FindByUserName(userName string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	// FindByIDWithRoles 预加载用户的角色信息
	FindByIDWithRoles(id uint) (*model.User, error)
	// FindPage 分页查询用户（预加载角色），支持按用户名模糊过滤
	FindPage(offset, limit int, userName string) ([]*model.User, int64, error)
}

// RoleRepository 角色数据访问接口
type RoleRepository interface {
	FindByCode(code string) (*model.Role, error)
	FindByUserID(userID uint) ([]*model.Role, error)
	Create(role *model.Role) error
	// FindPage 分页查询角色，支持按角色名模糊过滤
	FindPage(offset, limit int, roleName string) ([]*model.Role, int64, error)
}

// MenuRepository 菜单数据访问接口
type MenuRepository interface {
	// FindByRoleIDs 根据角色ID列表查询有权限的菜单（去重）
	FindByRoleIDs(roleIDs []uint) ([]*model.Menu, error)
	// FindAll 查询所有菜单
	FindAll() ([]*model.Menu, error)
	Create(menu *model.Menu) error
}

// PermissionRepository 权限数据访问接口
type PermissionRepository interface {
	// FindByRoleIDs 根据角色ID列表查询有权限的操作标识（去重）
	FindByRoleIDs(roleIDs []uint) ([]*model.Permission, error)
	// FindByRoleIDsAndMenuIDs 根据角色和菜单查询权限（用于构建菜单树的 authList）
	FindByRoleIDsAndMenuIDs(roleIDs []uint, menuIDs []uint) ([]*model.Permission, error)
	// FindByMenuID 查询某菜单下的所有权限
	FindByMenuID(menuID uint) ([]*model.Permission, error)
	Create(perm *model.Permission) error
}

// Repository 聚合所有数据访问接口，便于依赖注入
type Repository struct {
	User       UserRepository
	Role       RoleRepository
	Menu       MenuRepository
	Permission PermissionRepository
}

// NewRepository 创建 Repository 聚合实例
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:       NewUserRepository(db),
		Role:       NewRoleRepository(db),
		Menu:       NewMenuRepository(db),
		Permission: NewPermissionRepository(db),
	}
}
