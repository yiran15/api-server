package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/yiran15/api-server/model"
	"golang.org/x/oauth2"
)

var (
	oauthConfig = &oauth2.Config{
		ClientID:     "cli_a76d4e0cab38100c",
		ClientSecret: "Q04VFyNrQedJ3wrftlFIHfkDCWeisLWj",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.feishu.cn/open-apis/authen/v1/authorize",
			TokenURL: "https://open.feishu.cn/open-apis/authen/v2/oauth/token",
		},
		RedirectURL: "http://10.0.0.10:8080/oauth/callback",
	}
	state    = "random_state"
	verifier = "code_verifier"
)

type FeishuOauth struct {
	OAuthConfig *oauth2.Config
	State       string
	Verifier    string
}

func NewFeishuOauth() *FeishuOauth {
	return &FeishuOauth{
		OAuthConfig: oauthConfig,
		State:       state,
		Verifier:    verifier,
	}
}

func (f *FeishuOauth) GetAuthUrl() string {
	return f.OAuthConfig.AuthCodeURL(f.State, oauth2.S256ChallengeOption(f.Verifier))
}

func (f *FeishuOauth) ExchangeToken(ctx context.Context, state, code string) (*oauth2.Token, error) {
	if state != f.State {
		return nil, errors.New("state is not match")
	}
	if code == "" {
		return nil, errors.New("code is empty")
	}
	return f.OAuthConfig.Exchange(ctx, code, oauth2.VerifierOption(f.Verifier))
}

func (f *FeishuOauth) GetUserInfo(ctx context.Context, token *oauth2.Token) (*model.FeiShuUser, error) {
	client := f.OAuthConfig.Client(ctx, token)
	req, err := http.NewRequest("GET", "https://open.feishu.cn/open-apis/authen/v1/user_info", nil)
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
	var feishuUser model.FeiShuUser
	if err := json.Unmarshal(body, &feishuUser); err != nil {
		return nil, err
	}
	return &feishuUser, nil
}
