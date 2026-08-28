package repository

import (
	"golang-admin-best/internal/model"

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

func (r *roleRepository) FindPage(offset, limit int, roleName string) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64

	query := r.db.Model(&model.Role{})
	if roleName != "" {
		query = query.Where("role_name LIKE ?", "%"+roleName+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("id ASC").
		Offset(offset).Limit(limit).
		Find(&roles).Error
	return roles, total, err
}

func (r *roleRepository) Create(role *model.Role) error {
	return r.db.Create(role).Error
}
