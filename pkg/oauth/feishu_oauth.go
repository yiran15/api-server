package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/yiran15/api-server/base/conf"
	"github.com/yiran15/api-server/base/helper"
	"github.com/yiran15/api-server/model"
	"golang.org/x/oauth2"
)

type FeishuOauth struct {
	Name        string
	UserInfoUrl string
	State       string
	OAuthConfig *oauth2.Config
}

func NewFeishuOauth() (*FeishuOauth, error) {
	oauthConfig, err := conf.GetOauth2Config()
	if err != nil {
		return nil, err
	}
	state := conf.GetOauth2State()
	userInfoUrl, err := conf.GetOauth2UserInfoUrl()
	if err != nil {
		return nil, err
	}
	name := conf.GetOauth2Name()
	return &FeishuOauth{
		Name:        name,
		OAuthConfig: oauthConfig,
		UserInfoUrl: userInfoUrl,
		State:       state,
	}, nil
}

func (f *FeishuOauth) GetAuthUrl() string {
	return f.OAuthConfig.AuthCodeURL(f.State)
}

func (f *FeishuOauth) ExchangeToken(ctx context.Context, state, code string) (*oauth2.Token, error) {
	if state != f.State {
		return nil, errors.New("state is not match")
	}
	if code == "" {
		return nil, errors.New("code is empty")
	}
	return f.OAuthConfig.Exchange(ctx, code)
}

func (f *FeishuOauth) GetUserInfo(ctx context.Context, token *oauth2.Token) (*model.FeiShuUser, error) {
	client := f.OAuthConfig.Client(ctx, token)
	req, err := http.NewRequest("GET", f.UserInfoUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("body", string(body))

	var res helper.HttpResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return nil, errors.New(res.Msg)
	}

	feishuUser, err := helper.UnmarshalData[model.FeiShuUser](res.Data)
	if err != nil {
		return nil, err
	}
	return feishuUser, nil
}
