package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type DriverRedisCluster struct {
	client *redis.ClusterClient
}

func (d *DriverRedisCluster) SetUp(config Config) error {
	d.client = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    config.EndpointList,
		Password: config.Password,
	})

	_, err := d.client.Ping(d.client.Context()).Result()
	if err != nil {
		fmt.Printf("ping redis client err\n")
		return err
	}

	return nil
}

func (d *DriverRedisCluster) Get(key string) (data string, err error) {
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

func (d *DriverRedisCluster) Set(key string, value interface{}) error {
	return d.client.Set(ctx, key, value, time.Hour/2).Err()
}
