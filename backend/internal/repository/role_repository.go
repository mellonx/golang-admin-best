package repository

import (
	"art-design-pro-api/internal/model"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindByCode(code string) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("role_code = ?", code).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindByUserID(userID uint) ([]*model.Role, error) {
	var roles []*model.Role
	err := r.db.
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

func (r *roleRepository) Create(role *model.Role) error {
	return r.db.Create(role).Error
}
