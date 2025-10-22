package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

func Authenticate(ctx context.Context, authUrl, user, pass, realm, clientId, clientSecret string) (*http.Client, error) {
	grant_type := "client_credentials"
	if user != "" && pass != "" {
		grant_type = "password"
	}

	form := url.Values{}
	form.Add("grant_type", grant_type)
	form.Add("client_id", clientId)
	form.Add("client_secret", clientSecret)
	if user != "" && pass != "" {
		form.Add("username", user)
		form.Add("password", pass)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", authUrl, realm), strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		panic(err)
	}
	c := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	loginResponse := LoginResponse{}

	err = json.Unmarshal(respBytes, &loginResponse)
	if err != nil {
		panic(err)
	}

	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: loginResponse.AccessToken,
		TokenType:   loginResponse.TokenType,
	}))

	if loginResponse.AccessToken != "" {
		log.Println("Log in Successfully")
	} else {
		log.Printf("Log in failed - %s\n", string(respBytes))
		return nil, fmt.Errorf("login failed: %s", string(respBytes))
	}

	return client, nil
}

type LoginResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresInt       int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}
