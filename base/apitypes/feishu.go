package apitypes

import "github.com/yiran15/api-server/model"

type FeiShuLoginResponse struct {
	User  *model.FeiShuUser `json:"user"`
	Token string            `json:"token"`
}
