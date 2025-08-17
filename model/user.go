package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	ModelNameUser = "User"
	PreloadUsers  = "Users"
)

type User struct {
	ID         int64          `gorm:"column:id;primarykey;autoIncrement" json:"id"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
	Name       string         `gorm:"column:name;comment:用户名称;uniqueIndex;size:50" json:"name"`
	NickName   string         `gorm:"column:nick_name;comment:用户昵称;size:50" json:"nickName"`
	Department string         `gorm:"column:department;comment:用户部门;size:50" json:"department"`
	Email      string         `gorm:"column:email;comment:邮箱;uniqueIndex;size:100" json:"email"`
	Password   string         `gorm:"column:password;comment:用户密码;size:255" json:"-"`
	Avatar     string         `gorm:"column:avatar;comment:用户头像;size:1024" json:"avatar"`
	Mobile     string         `gorm:"column:mobile;comment:用户手机号;size:20" json:"mobile"`
	Status     *int           `gorm:"column:status;comment:用户状态,1可用,2禁用;size:1;default:1" json:"status"`
	Roles      []*Role        `gorm:"many2many:user_roles" json:"roles,omitempty"`
}

func (receiver *User) TableName() string {
	return "users"
}

type FeiShuUser struct {
	ID              int64          `gorm:"column:id;primarykey;autoIncrement" json:"id"`
	CreatedAt       time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt       time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
	AvatarBig       string         `gorm:"column:avatar_big;comment:飞书用户avatar_big;index" json:"avatar_big"`
	AvatarMiddle    string         `gorm:"column:avatar_middle;comment:飞书用户avatar_middle;index" json:"avatar_middle"`
	AvatarThumb     string         `gorm:"column:avatar_thumb;comment:飞书用户avatar_thumb;index" json:"avatar_thumb"`
	AvatarUrl       string         `gorm:"column:avatar_url;comment:飞书用户avatar_url;index" json:"avatar_url"`
	Email           string         `gorm:"column:email;comment:飞书用户email;index" json:"email"`
	EmployeeNo      string         `gorm:"column:employee_no;comment:飞书用户employee_no;index" json:"employee_no"`
	EnName          string         `gorm:"column:en_name;comment:飞书用户en_name;index" json:"en_name"`
	EnterpriseEmail string         `gorm:"column:enterprise_email;comment:飞书用户enterprise_email;index" json:"enterprise_email"`
	Mobile          string         `gorm:"column:mobile;comment:飞书用户mobile;index" json:"mobile"`
	Name            string         `gorm:"column:name;comment:飞书用户name;index" json:"name"`
	OpenID          string         `gorm:"column:open_id;comment:飞书用户open_id;index" json:"open_id"`
	TenantKey       string         `gorm:"column:tenant_key;comment:飞书用户tenant_key;index" json:"tenant_key"`
	UnionID         string         `gorm:"column:union_id;comment:飞书用户union_id;index" json:"union_id"`
	UserID          string         `gorm:"column:user_id;comment:飞书用户ID;index" json:"user_id"`
	UsersID         int64          `gorm:"column:users_id;comment:关联users表的用户id;index" json:"usersId"`
	User            *User          `gorm:"foreignKey:UsersID;references:ID" json:"user"`
	Status          *int           `gorm:"column:status;comment:用户状态,1可用,2禁用,3未激活;size:1;default:1" json:"status"`
}

func (receiver *FeiShuUser) TableName() string {
	return "feishu_users"
}
