package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
)

var (
	ctx         = context.Background()
	redisClient *redis.Client
)

func main() {
	// 1. 初始化 Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("無法連線到 Redis: %v", err)
	}

	// 2. 連線 NATS
	nc, err := nats.Connect(nats.DefaultURL) // 預設: nats://127.0.0.1:4222
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Drain()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	// 3. 建立 Stream
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "EVENTS",
		Subjects: []string{"events.*"},
		Storage:  nats.FileStorage,
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		log.Fatal(err)
	}

	// 4. 發布訊息（帶 event-id header）
	for i := 1; i <= 3; i++ {
		msgID := fmt.Sprintf("event-%d", i)
		msg := nats.NewMsg("events.created")
		msg.Data = []byte(fmt.Sprintf("payload-%d", i))
		msg.Header = nats.Header{"event-id": []string{msgID}}
		ack, err := js.PublishMsg(msg, nats.MsgId(msgID))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("[Publish] %s stored at seq %d\n", msgID, ack.Sequence)
	}

	// 5. 建立 Consumer (Push 訂閱)
	_, err = js.Subscribe("events.*", func(m *nats.Msg) {
		eventID := m.Header.Get("event-id")
		if eventID == "" {
			log.Printf("[Warning] 訊息缺少 event-id，略過")
			m.Ack()
			return
		}

		// 檢查是否已處理
		if alreadyProcessed(eventID) {
			fmt.Printf("[Duplicate] %s 已處理過，忽略\n", eventID)
			m.Ack()
			return
		}

		// 處理訊息
		fmt.Printf("[Process] %s: %s\n", eventID, string(m.Data))
		time.Sleep(500 * time.Millisecond) // 模擬業務邏輯

		// 標記已處理
		markProcessed(eventID, 10*time.Second)

		// Ack
		m.Ack()
	}, nats.Durable("EVENTS_CONSUMER"), nats.ManualAck())
	if err != nil {
		log.Fatal(err)
	}

	select {} // 阻塞
}

// 已處理檢查
func alreadyProcessed(eventID string) bool {
	val, err := redisClient.Get(ctx, "processed:"+eventID).Result()
	if err == redis.Nil {
		return false
	}
	if err != nil {
		log.Printf("Redis 錯誤: %v", err)
		return false
	}
	return val == "1"
}

// 標記已處理
func markProcessed(eventID string, ttl time.Duration) {
	err := redisClient.Set(ctx, "processed:"+eventID, "1", ttl).Err()
	if err != nil {
		log.Printf("Redis Set 錯誤: %v", err)
	}
}
