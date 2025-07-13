package model

import (
	"time"

	"gorm.io/gorm"
)

const PreloadRoles = "Roles"

type Role struct {
	ID          int64          `gorm:"column:id;primarykey" json:"id,omitempty"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"createdAt,omitempty"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updatedAt,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
	Name        string         `gorm:"column:name" json:"name,omitempty"`
	Description string         `gorm:"column:description" json:"description,omitempty"`
	Users       []*User        `gorm:"many2many:user_roles" json:"users,omitempty"`
	Apis        []*Api         `gorm:"many2many:role_apis" json:"apis,omitempty"`
}

func (receiver *Role) TableName() string {
	return "roles"
}

func (receiver *Role) AssociationModelName() string {
	return "Roles"
}
