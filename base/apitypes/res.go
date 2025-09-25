package apitypes

type Response struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg,omitempty"`
	Data      any    `json:"data,omitempty"`
	RequestId string `json:"requestId,omitempty"`
	Error     any    `json:"error,omitempty"`
}

func NewResponse(code int, msg, requestId string, data any, error any) *Response {
	return &Response{
		Code:      code,
		Msg:       msg,
		Data:      data,
		RequestId: requestId,
		Error:     error,
	}
}
