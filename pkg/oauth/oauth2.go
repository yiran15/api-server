package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/yiran15/api-server/base/conf"
	"github.com/yiran15/api-server/base/constant"
	"github.com/yiran15/api-server/base/helper"
	"github.com/yiran15/api-server/model"
	"golang.org/x/oauth2"
)

type OAuth2 struct {
	Name        string
	UserInfoUrl string
	OAuthConfig *oauth2.Config
}

func NewOAuth2() (*OAuth2, error) {
	if !conf.GetOauthEnable() {
		return nil, nil
	}
	oauthConfig, err := conf.GetOauth2Config()
	if err != nil {
		return nil, err
	}
	userInfoUrl, err := conf.GetOauth2UserInfoUrl()
	if err != nil {
		return nil, err
	}
	name := conf.GetOauth2Name()
	return &OAuth2{
		Name:        name,
		OAuthConfig: oauthConfig,
		UserInfoUrl: userInfoUrl,
	}, nil
}

func (f *OAuth2) Redirect(state string) string {
	return f.OAuthConfig.AuthCodeURL(state)
}

func (f *OAuth2) Auth(ctx context.Context, state, code string) (*oauth2.Token, error) {
	ctxState, ok := ctx.Value(constant.StateContextKey).(string)
	if !ok {
		return nil, errors.New("state not found")
	}
	if ctxState == "" {
		return nil, errors.New("state is empty")
	}
	if state != ctxState {
		return nil, errors.New("state is not match")
	}
	return f.OAuthConfig.Exchange(ctx, code)
}

func (f *OAuth2) UserInfo(ctx context.Context, token *oauth2.Token) (any, error) {
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
	var kcUser model.KeycloakUser
	if err := json.Unmarshal(body, &kcUser); err == nil && kcUser.Sub != "" {
		return &kcUser, nil
	}

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
