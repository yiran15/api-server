package apitypes

type IDRequest struct {
	ID int64 `uri:"id" validate:"required"`
}

type Pagination struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

type ListResponse struct {
	*Pagination
	Total int64 `json:"total"`
}
