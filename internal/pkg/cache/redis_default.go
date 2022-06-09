package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type DriverRedisDefault struct {
	client *redis.Client
}

func (d *DriverRedisDefault) SetUp(config Config) error {
	d.client = redis.NewClient(&redis.Options{
		Addr:     config.EndpointList[0],
		Password: config.Password,
	})

	_, err := d.client.Ping(d.client.Context()).Result()
	if err != nil {
		fmt.Printf("ping redis cluster client err\n")
		return err
	}

	return nil
}

func (d *DriverRedisDefault) Get(key string) (data string, err error) {
	cmd := d.client.Get(ctx, key)
	if cmd == nil {
		err = errors.New("redis client get err")
		return
	}
	data, err = cmd.Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
		}
		return
	}

	return
}

func (d *DriverRedisDefault) Set(key string, value interface{}) error {
	return d.client.Set(ctx, key, value, time.Hour/2).Err()
}
