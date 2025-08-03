package apitypes

import "github.com/yiran15/api-server/model"

type ApiCreateRequest struct {
	Name        string `json:"name" validate:"required"`
	Path        string `json:"path" validate:"required,uri"`
	Method      string `json:"method" validate:"required,oneof=GET POST PUT DELETE"`
	Description string `json:"description"`
}

type ApiUpdateRequest struct {
	*IDRequest
	Description string `json:"description"`
}

type ApiListRequest struct {
	*Pagination
	Name      string `form:"name"`
	Path      string `form:"path" validate:"omitempty,uri"`
	Method    string `form:"method" validate:"omitempty,oneof=GET POST PUT DELETE"`
	Sort      string `form:"sort" binding:"omitempty,oneof=id name path method created_at updated_at"`
	Direction string `form:"direction" binding:"omitempty,oneof=asc desc"`
}

type ApiListResponse struct {
	*ListResponse
	List []*model.Api `json:"list"`
}
