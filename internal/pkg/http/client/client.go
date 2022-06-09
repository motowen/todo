package client

import "github.com/go-resty/resty/v2"

var client *resty.Client

func Setup() error {
	client = resty.New()
	client.SetRetryCount(3)

	return nil
}

func NewHTTPRequest() *resty.Request {
	return client.R()
}

func Get() *resty.Client {
	return client
}
