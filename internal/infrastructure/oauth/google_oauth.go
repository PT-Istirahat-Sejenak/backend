package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type GooogleOauth struct {
	clientID     string
	clientSecret string
	redirectURL  string
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func NewGoogleOauth(clientID, clientSecret, redirectURL string) *GooogleOauth {
	return &GooogleOauth{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
	}
}

func (g *GooogleOauth) GetAuthURL() string {
	return fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=profile+email&access_type=offline",
		g.clientID,
		g.redirectURL,
	)
}

func (g *GooogleOauth) GetUserInfo(ctx context.Context, code string) (*GoogleUserInfo, error) {
	// menukar kode dengan token
	tokenURL := "https://oauth2.googleapis.com/token"
	values := url.Values{}
	values.Add("code", code)
	values.Add("client_id", g.clientID)
	values.Add("client_secret", g.clientSecret)
	values.Add("redirect_uri", g.redirectURL)
	values.Add("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to exchange code for token")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
	}

	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return nil, err
	}

	// Use access token to get user info
	userInfoURL := "https://www.googleapis.com/oauth2/v1/userinfo"
	req, err = http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+tokenResp.AccessToken)

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get user info")
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo GoogleUserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}
