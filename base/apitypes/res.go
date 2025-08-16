package apitypes

type Response struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg,omitempty"`
	Data      any    `json:"data,omitempty"`
	RequestId string `json:"requestId,omitempty"`
}

func NewResponse(code int, msg, requestId string, data any) *Response {
	return &Response{
		Code:      code,
		Msg:       msg,
		Data:      data,
		RequestId: requestId,
	}
}
