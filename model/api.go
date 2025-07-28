package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	PreloadApis = "Apis"
)

type Api struct {
	ID          int64          `gorm:"column:id;primarykey" json:"id,omitempty"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"createdAt,omitempty"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updatedAt,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
	Name        string         `gorm:"column:name" json:"name,omitempty"`
	Path        string         `gorm:"column:path" json:"path,omitempty"`
	Method      string         `gorm:"column:method" json:"method,omitempty"`
	Description string         `gorm:"column:description" json:"description,omitempty"`
	Roles       []*Role        `gorm:"many2many:role_apis" json:"roles,omitempty"`
}

func (*Api) TableName() string {
	return "apis"
}
