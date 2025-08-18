package model

import (
	"time"

	"gorm.io/gorm"
)

var (
	FeiShuUserStatusActive   = 1
	FeiShuUserStatusDisabled = 2
	FeiShuUserStatusInactive = 3
)

type FeiShuUser struct {
	UID             int64          `gorm:"column:uid;primarykey;comment:关联users表中的用户id" json:"uid"`
	User            *User          `gorm:"foreignKey:UID;references:ID" json:"user"`
	CreatedAt       time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt       time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
	AvatarBig       string         `gorm:"column:avatar_big;comment:飞书用户avatar_big" json:"avatar_big"`
	AvatarMiddle    string         `gorm:"column:avatar_middle;comment:飞书用户avatar_middle" json:"avatar_middle"`
	AvatarThumb     string         `gorm:"column:avatar_thumb;comment:飞书用户avatar_thumb" json:"avatar_thumb"`
	AvatarUrl       string         `gorm:"column:avatar_url;comment:飞书用户avatar_url" json:"avatar_url"`
	Email           string         `gorm:"column:email;comment:飞书用户email" json:"email"`
	EmployeeNo      string         `gorm:"column:employee_no;comment:飞书用户employee_no" json:"employee_no"`
	EnName          string         `gorm:"column:en_name;comment:飞书用户en_name" json:"en_name"`
	EnterpriseEmail string         `gorm:"column:enterprise_email;comment:飞书用户enterprise_email" json:"enterprise_email"`
	Mobile          string         `gorm:"column:mobile;comment:飞书用户mobile" json:"mobile"`
	Name            string         `gorm:"column:name;comment:飞书用户name" json:"name"`
	OpenID          string         `gorm:"column:open_id;comment:飞书用户open_id" json:"open_id"`
	TenantKey       string         `gorm:"column:tenant_key;comment:飞书用户tenant_key" json:"tenant_key"`
	UnionID         string         `gorm:"column:union_id;comment:飞书用户union_id" json:"union_id"`
	UserID          string         `gorm:"column:user_id;comment:飞书用户ID;index:idx_user_id_status,priority:1" json:"user_id"`
	Status          *int           `gorm:"column:status;comment:用户状态,1可用,2禁用,3未激活;size:1;default:1;index:idx_user_id_status,priority:2;" json:"status"`
}

func (receiver *FeiShuUser) TableName() string {
	return "feishu_users"
}
