package apitypes

type IDRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type Pagination struct {
	Page     int `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize int `form:"pageSize" json:"pageSize" binding:"omitempty,min=1,max=100"`
}

type ListResponse struct {
	*Pagination
	Total int64 `json:"total"`
}
