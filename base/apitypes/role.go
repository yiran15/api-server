package apitypes

import "github.com/yiran15/api-server/model"

type RoleCreateRequest struct {
	Name        string  `json:"name" validate:"required,ascii"`
	Description string  `json:"description"`
	Apis        []int64 `json:"apis"`
}

type RoleUpdateRequest struct {
	*IDRequest
	Description string  `json:"description"`
	Apis        []int64 `json:"apis"`
}

type RoleListRequest struct {
	*Pagination
	Name      string `form:"name"`
	Sort      string `form:"sort" binding:"omitempty,oneof=id name created_at updated_at"`
	Direction string `form:"direction" binding:"omitempty,oneof=asc desc"`
}

type RoleListResponse struct {
	*ListResponse
	List []*model.Role `json:"list"`
}
