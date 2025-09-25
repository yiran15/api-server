package apitypes

import (
	"github.com/yiran15/api-server/model"
)

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserLoginResponse struct {
	User  *model.User `json:"user"`
	Token string      `json:"token"`
}

type UserCreateRequest struct {
	Name     string   `json:"name" binding:"required"`
	NickName string   `json:"nickName"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=8"`
	Avatar   string   `json:"avatar"`
	Mobile   string   `json:"mobile" binding:"omitempty,mobile"`
	RolesID  *[]int64 `json:"rolesID"`
}

type UserUpdateAdminRequest struct {
	ID int64 `uri:"id" binding:"required"`
	*UserUpdateSelfRequest
	Status  int      `json:"status" binding:"omitempty,oneof=1 2"`
	RolesID *[]int64 `json:"rolesID" binding:"omitempty"`
}

type UserUpdateSelfRequest struct {
	Name        string `json:"name"`
	NickName    string `json:"nickName"`
	Email       string `json:"email" binding:"omitempty,email"`
	OldPassword string `json:"oldPassword" binding:"omitempty,min=8"`
	Password    string `json:"password" binding:"omitempty,min=8"`
	Avatar      string `json:"avatar"`
	Mobile      string `json:"mobile" binding:"omitempty,mobile"`
}

type UserUpdateStatusRequest struct {
	ID     int64 `uri:"id" binding:"required"`
	Status int   `json:"status" binding:"required,oneof=1 2"`
}

type UserIdRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type UserListRequest struct {
	*Pagination
	Name       string `form:"name" binding:"user_list"`
	Email      string `form:"email" binding:"omitempty,email"`
	Mobile     string `form:"mobile" binding:"omitempty,mobile"`
	Department string `form:"department"`
	Sort       string `form:"sort" binding:"omitempty,oneof=id name created_at updated_at nick_name email mobile"`
	Direction  string `form:"direction" binding:"omitempty,oneof=asc desc"`
	Status     int    `form:"status" binding:"omitempty,oneof=0 1 2"`
}

type UserListResponse struct {
	*ListResponse
	List []*model.User `json:"list"`
}

type UserUpdateRoleRequest struct {
	ID      int64   `uri:"id" binding:"required"`
	RolesID []int64 `json:"rolesID" binding:"required"`
}

type OAuthLoginRequest struct {
	Code string `form:"code" binding:"required"`
}

type OauthLoginResponse struct {
	User  any    `json:"user"`
	Token string `json:"token"`
}

type OAuthActivateRequest struct {
	ID              string `uri:"id" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=8"`
}
