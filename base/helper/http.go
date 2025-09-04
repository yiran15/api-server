package helper

import "encoding/json"

// 请求 http 服务返回的通用数据结构
type HttpResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func UnmarshalData[T any](data json.RawMessage) (*T, error) {
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
