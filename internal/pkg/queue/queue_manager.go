package queue

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
)

var (
	instance *Manager
)

func GetInstance() *Manager {
	return instance
}

type Manager struct {
	connect *nats.Conn
}

type Config struct {
	Url string
}

func (manager *Manager) Setup(config Config) error {

	nc, err := nats.Connect(config.Url)
	if err != nil {
		fmt.Printf("queue connect fail, %+v\n", err)
		return err
	}

	instance = &Manager{
		connect: nc,
	}

	return nil
}

func (manager *Manager) Publish(sub string, msg interface{}) error {

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = manager.connect.Publish(sub, data)
	if err != nil {
		return err
	}

	return nil
}

func (manager *Manager) Subscribe(sub string, msgHandler func(sub string, msg []byte)) (*nats.Subscription, error) {

	subscription, err := manager.connect.Subscribe(sub, func(msg *nats.Msg) {
		msgHandler(msg.Subject, msg.Data)
	})
	if err != nil {
		return nil, err
	}

	return subscription, nil
}
