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
	Name string `form:"name"`
	*SortParam
}

type RoleListResponse struct {
	*ListResponse
	List []*model.Role `json:"list"`
}
