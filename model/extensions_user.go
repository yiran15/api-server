package model

import (
	"time"

	"gorm.io/gorm"
)

type FeiShuUser struct {
	UID             int64          `gorm:"column:uid;primarykey;comment:关联users表中的用户id" json:"uid"`
	User            *User          `gorm:"foreignKey:UID;references:ID" json:"user"`
	CreatedAt       time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt       time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
	AvatarBig       string         `gorm:"column:avatar_big;size:255;comment:飞书用户avatar_big" json:"avatar_big"`
	AvatarMiddle    string         `gorm:"column:avatar_middle;size:255;comment:飞书用户avatar_middle" json:"avatar_middle"`
	AvatarThumb     string         `gorm:"column:avatar_thumb;size:255;comment:飞书用户avatar_thumb" json:"avatar_thumb"`
	AvatarUrl       string         `gorm:"column:avatar_url;size:255;comment:飞书用户avatar_url" json:"avatar_url"`
	Email           string         `gorm:"column:email;size:255;comment:飞书用户email" json:"email"`
	EmployeeNo      string         `gorm:"column:employee_no;size:255;comment:飞书用户employee_no" json:"employee_no"`
	EnName          string         `gorm:"column:en_name;size:255;comment:飞书用户en_name" json:"en_name"`
	EnterpriseEmail string         `gorm:"column:enterprise_email;size:255;comment:飞书用户enterprise_email" json:"enterprise_email"`
	Mobile          string         `gorm:"column:mobile;size:255;comment:飞书用户mobile" json:"mobile"`
	Name            string         `gorm:"column:name;size:255;comment:飞书用户name" json:"name"`
	OpenID          string         `gorm:"column:open_id;size:255;comment:飞书用户open_id" json:"open_id"`
	TenantKey       string         `gorm:"column:tenant_key;size:255;comment:飞书用户tenant_key" json:"tenant_key"`
	UnionID         string         `gorm:"column:union_id;size:255;comment:飞书用户union_id" json:"union_id"`
	UserID          string         `gorm:"column:user_id;size:255;comment:飞书用户ID;index:idx_user_id_status,priority:1" json:"user_id"`
}

func (receiver *FeiShuUser) TableName() string {
	return "feishu_users"
}

type KeycloakUser struct {
	Sub               string   `json:"sub"`
	EmailVerified     bool     `json:"email_verified"`
	Roles             []string `json:"roles"`
	Name              string   `json:"name"`
	PreferredUsername string   `json:"preferred_username"`
	GivenName         string   `json:"given_name"`
	FamilyName        string   `json:"family_name"`
	Email             string   `json:"email"`
	Group             []string `json:"group"`
}
