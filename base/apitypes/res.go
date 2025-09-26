package apitypes

type Response struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg,omitempty"`
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

type ResponseOptions func(*Response)

func WithMsg(msg string) ResponseOptions {
	return func(r *Response) {
		r.Msg = msg
	}
}

func WithData(data any) ResponseOptions {
	return func(r *Response) {
		r.Data = data
	}
}

func WithError(err string) ResponseOptions {
	return func(r *Response) {
		r.Error = err
	}
}

func NewResponseWithOpts(code int, opts ...ResponseOptions) *Response {
	res := &Response{
		Code: code,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}
