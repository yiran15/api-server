package oauth_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestMain(t *testing.T) {
	ctx := context.Background()

	conf := &oauth2.Config{
		ClientID:     "cli_a76d4e0cab38100c",
		ClientSecret: "Q04VFyNrQedJ3wrftlFIHfkDCWeisLWj",
		// Scopes:       []string{"all"},
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://open.feishu.cn/open-apis/authen/v2/oauth/token",
			AuthURL:  "https://accounts.feishu.cn/open-apis/authen/v1/authorize",
		},
		RedirectURL: "http://localhost:8080/callback",
	}

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)
	var (
		resp *http.Response
		err  error
	)
	if resp, err = http.Get(url); err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}
	fmt.Println(data["code"])

	httpClient := &http.Client{Timeout: 2 * time.Second}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	tok, err := conf.Exchange(ctx, data["code"].(string), oauth2.SetAuthURLParam("state", "state"))
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(ctx, tok)
	fmt.Println(client)
}
