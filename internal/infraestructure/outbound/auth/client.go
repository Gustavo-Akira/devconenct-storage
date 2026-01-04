package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type IAuthClient interface {
	GetProfile(token string) (*int64, error)
}

type AuthClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAuthClient(baseURL string) *AuthClient {
	client := http.DefaultClient
	return &AuthClient{
		baseURL:    baseURL,
		httpClient: client,
	}
}

type ProfileMeResponse struct {
	ID int64 `json:"id"`
}

func (ac *AuthClient) GetProfile(token string) (*int64, error) {
	req, err := http.NewRequest("GET", ac.baseURL+"/v1/dev-profiles/profile", nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{
		Name:  "jwt",
		Value: token,
	})

	return doAuthRequest(ac.httpClient, req)
}

func doAuthRequest(client *http.Client, req *http.Request) (*int64, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var authProfile ProfileMeResponse
	if err := json.NewDecoder(resp.Body).Decode(&authProfile); err != nil {
		return nil, err
	}
	return &authProfile.ID, nil
}
