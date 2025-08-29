package apitypes

import (
	"github.com/yiran15/api-server/model"
)

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserLoginResponse struct {
	User  *model.User `json:"user"`
	Token string      `json:"token"`
}

type UserCreateRequest struct {
	Name     string   `json:"name" validate:"required"`
	NickName string   `json:"nickName"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=8"`
	Avatar   string   `json:"avatar"`
	Mobile   string   `json:"mobile" validate:"omitempty,mobile"`
	RolesID  *[]int64 `json:"rolesID"`
}

type UserUpdateAdminRequest struct {
	ID int64 `uri:"id" validate:"required"`
	*UserUpdateSelfRequest
	Status  int      `json:"status" validate:"omitempty,oneof=1 2"`
	RolesID *[]int64 `json:"rolesID" validate:"omitempty"`
}

type UserUpdateSelfRequest struct {
	Name        string `json:"name"`
	NickName    string `json:"nickName"`
	Email       string `json:"email" validate:"omitempty,email"`
	OldPassword string `json:"oldPassword" validate:"omitempty,min=8"`
	Password    string `json:"password" validate:"omitempty,min=8"`
	Avatar      string `json:"avatar"`
	Mobile      string `json:"mobile" validate:"omitempty,mobile"`
}

type UserUpdateStatusRequest struct {
	ID     int64 `uri:"id" validate:"required"`
	Status int   `json:"status" validate:"required,oneof=1 2"`
}

type UserIdRequest struct {
	ID int64 `uri:"id" validate:"required"`
}

type UserListRequest struct {
	*Pagination
	Name       string `form:"name" validate:"user_list"`
	Email      string `form:"email" validate:"omitempty,email"`
	Mobile     string `form:"mobile" validate:"omitempty,mobile"`
	Department string `form:"department"`
	Sort       string `form:"sort" binding:"omitempty,oneof=id name created_at updated_at nick_name email mobile"`
	Direction  string `form:"direction" binding:"omitempty,oneof=asc desc"`
	Status     int    `form:"status" validate:"omitempty,oneof=0 1 2"`
}

type UserListResponse struct {
	*ListResponse
	List []*model.User `json:"list"`
}

type UserUpdateRoleRequest struct {
	ID      int64   `uri:"id" validate:"required"`
	RolesID []int64 `json:"rolesID" validate:"required"`
}

type OAuthLoginRequest struct {
	Code  string `form:"code" validate:"required"`
	State string `form:"state" validate:"required"`
}

type OauthLoginResponse struct {
	User  any    `json:"user"`
	Token string `json:"token"`
}
