package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	ModelNameUser = "User"
	PreloadUsers  = "Users"
)

var (
	UserStatusActive   = 1
	UserStatusDisabled = 2
)

type User struct {
	ID         int64          `gorm:"column:id;primarykey;autoIncrement" json:"id"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
	Name       string         `gorm:"column:name;comment:用户名称;size:50" json:"name"`
	NickName   string         `gorm:"column:nick_name;comment:用户昵称;size:50" json:"nickName"`
	Department string         `gorm:"column:department;comment:用户部门;size:50" json:"department"`
	Email      string         `gorm:"column:email;comment:邮箱;size:100" json:"email"`
	Password   string         `gorm:"column:password;comment:用户密码;size:255" json:"-"`
	Avatar     string         `gorm:"column:avatar;comment:用户头像;size:1024" json:"avatar"`
	Mobile     string         `gorm:"column:mobile;comment:用户手机号;size:20" json:"mobile"`
	Status     *int           `gorm:"column:status;comment:用户状态,1可用,2禁用;size:1;default:1" json:"status"`
	Roles      []*Role        `gorm:"many2many:user_roles" json:"roles,omitempty"`
}

func (receiver *User) TableName() string {
	return "users"
}
