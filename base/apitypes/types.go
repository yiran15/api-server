package apitypes

type IDRequest struct {
	ID int64 `uri:"id" validate:"required"`
}

type Pagination struct {
	Page     int `form:"page" validate:"required"`
	PageSize int `form:"pageSize" validate:"required"`
}

type ListResponse struct {
	*Pagination
	Total int64 `json:"total"`
}

type SortParam struct {
	Sort      string `form:"sort" binding:"omitempty,oneof=id name created_at updated_at nick_name email mobile"` // 排序字段名，如 "createdAt", "name"
	Direction string `form:"direction" binding:"omitempty,oneof=asc desc"`                                        // 排序方向，"asc" 或 "desc"
}

func (receiver *SortParam) GetRoleSortField(fieldNum uint) string {
	switch fieldNum {
	case 1:
		return "id"
	case 2:
		return "created_at"
	case 3:
		return "updated_at"
	case 4:
		return "name"
	default:
		return "id"
	}
}

func (receiver *SortParam) GetApiSortField(fieldNum uint) string {
	switch fieldNum {
	case 1:
		return "id"
	case 2:
		return "created_at"
	case 3:
		return "updated_at"
	case 4:
		return "name"
	case 5:
		return "path"
	case 6:
		return "method"
	default:
		return "id"
	}
}
