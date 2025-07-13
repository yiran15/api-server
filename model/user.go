package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           int64          `gorm:"column:id;primarykey" json:"id"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
	Name         string         `gorm:"column:name;comment:用户名称;uniqueIndex;size:50" json:"name"`
	NickName     string         `gorm:"column:nick_name;comment:用户昵称;size:50" json:"nickName"`
	OpenId       string         `gorm:"column:open_id;comment:飞书用户ID;size:50" json:"openId"`
	Department   string         `gorm:"column:department;comment:部门;size:255" json:"department"`
	DepartmentId string         `gorm:"column:department_id;comment:部门ID;size:255" json:"departmentId"`
	Email        string         `gorm:"column:email;comment:邮箱;uniqueIndex;size:100" json:"email"`
	Password     string         `gorm:"column:password;comment:用户密码;size:255" json:"-"`
	Avatar       string         `gorm:"column:avatar;comment:用户头像;size:1024" json:"avatar"`
	Mobile       string         `gorm:"column:mobile;comment:用户手机号;size:20" json:"mobile"`
	Status       *int           `gorm:"column:status;comment:用户状态,1可用,2删除;size:1;default:1" json:"status"`
	Roles        []*Role        `gorm:"many2many:users_roles" json:"roles,omitempty"`
}

func (receiver *User) TableName() string {
	return "users"
}
