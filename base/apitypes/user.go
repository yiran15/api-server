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
	Name     string `json:"name" validate:"required"`
	NickName string `json:"nickName"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Avatar   string `json:"avatar"`
	Mobile   string `json:"mobile" validate:"omitempty,mobile"`
}

type UserUpdateRequest struct {
	ID       int64  `uri:"id" validate:"required"`
	Name     string `json:"name"`
	NickName string `json:"nickName"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=8"`
	Avatar   string `json:"avatar"`
	Mobile   string `json:"mobile" validate:"omitempty,mobile"`
}

type UserIdRequest struct {
	ID int64 `uri:"id" validate:"required"`
}

type UserListRequest struct {
	*Pagination
	Name   string `form:"name" validate:"user_list"`
	Email  string `form:"email" validate:"omitempty,email"`
	Mobile string `form:"mobile" validate:"omitempty,mobile"`
	*SortParam
}

type UserListResponse struct {
	*ListResponse
	List []*model.User `json:"list"`
}

type UserUpdateRoleRequest struct {
	ID      int64   `uri:"id" validate:"required"`
	RoleIds []int64 `json:"roleIds" validate:"required"`
}
