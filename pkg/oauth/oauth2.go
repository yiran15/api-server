package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"
	"github.com/yiran15/api-server/base/helper"
	"github.com/yiran15/api-server/model"
	"golang.org/x/oauth2"
)

type OAuth2ProviderConfig struct {
	UserInfoUrl  string   `mapstructure:"userInfoUrl"`
	ClientId     string   `mapstructure:"clientId"`
	ClientSecret string   `mapstructure:"clientSecret"`
	Scopes       []string `mapstructure:"scopes"`
	AuthUrl      string   `mapstructure:"authUrl"`
	TokenUrl     string   `mapstructure:"tokenUrl"`
	RedirectUrl  string   `mapstructure:"redirectUrl"`
}

type OAuth2 struct {
	Enable    bool
	Providers map[string]*Provider
}

type Provider struct {
	UserInfoUrl string
	OAuthConfig *oauth2.Config
}

func NewOAuth2() (*OAuth2, error) {
	enable := viper.GetBool("oauth2.enable")
	if !enable {
		return nil, nil
	}

	providerConfigs := make(map[string]*OAuth2ProviderConfig)
	if err := viper.UnmarshalKey("oauth2.providers", &providerConfigs); err != nil {
		return nil, fmt.Errorf("unmarshal oauth2.providers faild. err: %w", err)
	}

	providers := make(map[string]*Provider)
	for name, providerConfig := range providerConfigs {
		providers[name] = &Provider{
			UserInfoUrl: providerConfig.UserInfoUrl,
			OAuthConfig: &oauth2.Config{
				ClientID:     providerConfig.ClientId,
				ClientSecret: providerConfig.ClientSecret,
				Endpoint: oauth2.Endpoint{
					AuthURL:  providerConfig.AuthUrl,
					TokenURL: providerConfig.TokenUrl,
				},
				RedirectURL: providerConfig.RedirectUrl,
				Scopes:      providerConfig.Scopes,
			},
		}
	}

	return &OAuth2{Enable: enable, Providers: providers}, nil
}

func (f *OAuth2) Redirect(state string, provider string) string {
	p, ok := f.Providers[provider]
	if !ok {
		return ""
	}
	return p.OAuthConfig.AuthCodeURL(state)
}

func (f *OAuth2) Auth(ctx context.Context, code, provider string) (*oauth2.Token, error) {
	p, ok := f.Providers[provider]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", provider)
	}
	return p.OAuthConfig.Exchange(ctx, code)
}

func (f *OAuth2) UserInfo(ctx context.Context, token *oauth2.Token, provider string) (any, error) {
	p, ok := f.Providers[provider]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", provider)
	}
	client := p.OAuthConfig.Client(ctx, token)
	req, err := http.NewRequest("GET", p.UserInfoUrl, nil)
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
