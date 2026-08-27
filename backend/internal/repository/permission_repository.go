package repository

import (
	"art-design-pro-api/internal/model"

	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) FindByRoleIDs(roleIDs []uint) ([]*model.Permission, error) {
	var perms []*model.Permission
	err := r.db.
		Distinct().
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id IN ?", roleIDs).
		Find(&perms).Error
	return perms, err
}

func (r *permissionRepository) FindByRoleIDsAndMenuIDs(roleIDs []uint, menuIDs []uint) ([]*model.Permission, error) {
	var perms []*model.Permission
	err := r.db.
		Distinct().
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id IN ? AND permissions.menu_id IN ?", roleIDs, menuIDs).
		Find(&perms).Error
	return perms, err
}

func (r *permissionRepository) FindByMenuID(menuID uint) ([]*model.Permission, error) {
	var perms []*model.Permission
	err := r.db.Where("menu_id = ?", menuID).Find(&perms).Error
	return perms, err
}

func (r *permissionRepository) Create(perm *model.Permission) error {
	return r.db.Create(perm).Error
}
