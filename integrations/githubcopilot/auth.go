package githubcopilot

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"net/http"
	"time"
)

const (
	clientID         = "Iv1.b507a08c87ecfe98"
	deviceCodeURL    = "https://github.com/login/device/code"
	accessTokenURL   = "https://github.com/login/oauth/access_token"
	copilotApiKeyURL = "https://api.github.com/copilot_internal/v2/token"
)

type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type AccessTokenResponse struct {
	AccessToken      string `json:"access_token"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type CopilotTokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

// StartDeviceAuth initiates device flow
func StartDeviceAuth() (*DeviceCodeResponse, error) {
	payload := map[string]string{
		"client_id": clientID,
		"scope":     "read:user",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", deviceCodeURL, bytes.NewReader(body))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "CodeForgeAI.go/0.1")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result DeviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PollForAccessToken polls for access token
func PollForAccessToken(deviceCode string) (*AccessTokenResponse, error) {
	payload := map[string]string{
		"client_id":   clientID,
		"device_code": deviceCode,
		"grant_type":  "urn:ietf:params:oauth:grant-type:device_code",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", accessTokenURL, bytes.NewReader(body))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "CodeForgeAI.go/0.1")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result AccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCopilotToken exchanges refresh token for Copilot API key
func GetCopilotToken(refreshToken string) (*CopilotTokenResponse, error) {
	req, _ := http.NewRequest("GET", copilotApiKeyURL, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+refreshToken)
	req.Header.Set("User-Agent", "CodeForgeAI.go/0.1")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result CopilotTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
