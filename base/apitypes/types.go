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
	Sort      uint   `form:"sort"`                                         // 排序字段名，如 "createdAt", "name"
	Direction string `form:"direction" binding:"omitempty,oneof=ASC DESC"` // 排序方向，"asc" 或 "desc"
}

// GetSortField 根据字段序号返回排序字段名
// 1 id
// 2 created_at
// 3 updated_at
// 4 name
// 5 nick_name
// 6 email
// 7 mobile
func (receiver *SortParam) GetUserSortField(fieldNum uint) string {
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
		return "nick_name"
	case 6:
		return "email"
	case 7:
		return "mobile"
	default:
		return "id"
	}
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
