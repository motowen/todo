package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/config"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/http/client"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/model"
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

	req := client.NewHTTPRequest().
		SetBody(body).
		SetHeader("Content-Type", "application/x-www-form-urlencoded")

	resp, err := req.Post(api)
	if err != nil {
		return "", model.ServiceError.InternalServiceError(err.Error())
	}

	if resp.StatusCode() != http.StatusOK {
		return "", model.ServiceError.InternalServiceError(err.Error())
	}

	auth := GetAuthTokenResponse{}
	if err := json.Unmarshal(resp.Body(), &auth); err != nil {
		return "", model.ServiceError.InternalServiceError(err.Error())
	}

	authCache.setToken(auth.AccessToken)
	return auth.AccessToken, model.ServiceError.OK
}
