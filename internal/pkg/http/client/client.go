package client

import (
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	client *resty.Client
	once   sync.Once
)

// 初始化 Resty Client
func Setup() *resty.Client {
	once.Do(func() {
		client = resty.New()

		// 設定全域 Timeout，避免 API 卡住
		client.SetTimeout(5 * time.Second)

		// Retry 設定：2 次，500ms 間隔
		client.SetRetryCount(2)
		client.SetRetryWaitTime(500 * time.Millisecond)

		// Retry 條件：網路錯誤或 HTTP 5xx 424
		client.AddRetryCondition(func(r *resty.Response, err error) bool {
			if err != nil {
				return true // 網路層錯誤，重試
			}
			return r.StatusCode() >= http.StatusFailedDependency
		})
	})
	return client
}

// 取得新的 Request，確保 client 一定初始化
func NewRequest() *resty.Request {
	if client == nil {
		Setup()
	}
	return client.R()
}

// 取得全域 Client
func Get() *resty.Client {
	if client == nil {
		Setup()
	}
	return client
}
