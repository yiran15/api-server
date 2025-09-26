package apitypes

import "github.com/yiran15/api-server/model"

type ApiCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Path        string `json:"path" binding:"required,uri"`
	Method      string `json:"method" binding:"required,oneof=GET POST PUT DELETE *"`
	Description string `json:"description"`
}

type ApiUpdateRequest struct {
	*IDRequest
	Description string `json:"description"`
}

type ApiListRequest struct {
	*Pagination
	Name      string `form:"name"`
	Path      string `form:"path" binding:"omitempty,uri"`
	Method    string `form:"method" binding:"omitempty,oneof=GET POST PUT DELETE"`
	Sort      string `form:"sort" binding:"omitempty,oneof=id name path method created_at updated_at"`
	Direction string `form:"direction" binding:"omitempty,oneof=asc desc"`
}

type ApiListResponse struct {
	*ListResponse
	List []*model.Api `json:"list"`
}

type ServerApiData struct {
	ApiType []string             `json:"apiType"`
	ApiInfo map[string][]ApiInfo `json:"apiInfo"`
}

type ApiInfo struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Handler string `json:"handler"`
}
