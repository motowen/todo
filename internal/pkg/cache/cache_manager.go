package cache

import (
	"context"
	"net/url"
)

var (
	instance *Manager
	ctx      = context.Background()
)

type Manager struct {
	Driver
}

type Config struct {
	Type         string
	EndpointList []string
	Password     string
}

type Driver interface {
	SetUp(config Config) error
	Get(key string) (string, error)
	Set(key string, value interface{}) error
}

func GetInstance() *Manager {
	if instance == nil {
		instance = &Manager{}
	}
	return instance
}

func (manager *Manager) Setup(config Config) error {
	var driver Driver

	switch config.Type {
	case "redis_default":
		driver = &DriverRedisDefault{}
	case "redis_cluster":
		driver = &DriverRedisCluster{}
	}

	err := driver.SetUp(config)
	manager.Driver = driver

	return err
}

type SessionInfo struct {
	Id         string
	UserName   string
	IsLoggedIn bool
	FormData   url.Values
}

type ThirdPartyLoginInfo struct {
	Id    string
	State string
}
