package externalAccount

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go-base/internal/pkg/config"
	"go-base/internal/pkg/http/client"
	"go-base/internal/pkg/logger"
	"go-base/internal/pkg/model"
)

type AuthCache struct {
	IsSet          bool
	Token          string
	StartTimeInSec int64
}

var authCache AuthCache

func (a *AuthCache) getToken() (string, error) {
	if a.IsSet && a.StartTimeInSec+config.Env.AuthServiceCacheTTL > time.Now().Unix() {
		return a.Token, nil
	}

	return "", errors.New("no token")
}

func (a *AuthCache) setToken(token string) {
	a.IsSet = true
	a.Token = token
	a.StartTimeInSec = time.Now().Unix()
}

func (a *AuthCache) clear() {
	a.IsSet = false
	a.Token = ""
	a.StartTimeInSec = 0
}

type GetAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func GetAuthToken() (string, model.ServiceResp) {
	token, err := authCache.getToken()
	if err == nil {
		return token, model.ServiceError.OK
	}

	api := fmt.Sprintf("%s%s", config.Env.AuthServiceHost, "/$SS$/Services/OAuth/Token")
	body := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", config.Env.AuthClientID, config.Env.AuthClientSecret)

	resp, err := client.NewRequest().
		SetBody(body).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		Post(api)
	if err != nil {
		logger.Error.Printf("Failed to GetAuthToken: %v", err)
		return "", model.ServiceError.FailedDependencyError(model.ExternalGetAuthTokenFail)
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Error.Printf("GetAuthToken StatusCode: %v", resp.StatusCode())
		return "", model.ServiceError.FailedDependencyError(model.ExternalGetAuthTokenFail)
	}

	auth := GetAuthTokenResponse{}
	if err := json.Unmarshal(resp.Body(), &auth); err != nil {
		logger.Error.Printf("GetAuthToken Unmarshal: %v", err)
		return "", model.ServiceError.FailedDependencyError(model.ExternalGetAuthTokenParseFail)
	}

	// 檢查解析結果是否有效 - 這是關鍵的驗證步驟
	if auth.AccessToken == "" {
		logger.Error.Printf("GetAuthToken response missing access_token: %v, body: %s", auth, string(resp.Body()))
		return "", model.ServiceError.FailedDependencyError(model.ExternalGetAuthTokenParseFail)
	}

	logger.Info.Printf("GetAuthToken Success: %v", auth)
	authCache.setToken(auth.AccessToken)
	return auth.AccessToken, model.ServiceError.OK
}

// ClearAuthCache clears the authentication token cache (for testing)
func ClearAuthCache() {
	authCache.clear()
}
