package repository

import (
	"art-design-pro-api/internal/model"

	"gorm.io/gorm"
)

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) FindByRoleIDs(roleIDs []uint) ([]*model.Menu, error) {
	var menus []*model.Menu
	err := r.db.
		Distinct().
		Joins("JOIN role_menus ON role_menus.menu_id = menus.id").
		Where("role_menus.role_id IN ? AND menus.status = ?", roleIDs, 1).
		Order("menus.sort ASC, menus.id ASC").
		Find(&menus).Error
	return menus, err
}

func (r *menuRepository) FindAll() ([]*model.Menu, error) {
	var menus []*model.Menu
	err := r.db.Where("status = ?", 1).Order("sort ASC, id ASC").Find(&menus).Error
	return menus, err
}

func (r *menuRepository) Create(menu *model.Menu) error {
	return r.db.Create(menu).Error
}
