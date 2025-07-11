package model

import "time"

type Api struct {
	Id          int64     `gorm:"column:id;primarykey"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
	Name        string    `gorm:"column:name"`
	Path        string    `gorm:"column:path"`
	Method      string    `gorm:"column:method"`
	Description string    `gorm:"column:description"`
}

func (*Api) TableName() string {
	return "apis"
}
