package model

import "time"

// User 用户表
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserName  string    `gorm:"uniqueIndex;size:50;not null" json:"userName"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Email     string    `gorm:"size:100" json:"email"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	Status    int8      `gorm:"default:1;comment:1:启用 0:禁用" json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Roles []*Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

func (User) TableName() string {
	return "users"
}
