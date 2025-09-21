package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// 1. 建立 NATS 連線
	nc, err := nats.Connect(nats.DefaultURL) // 預設: nats://127.0.0.1:4222
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Drain()

	// 2. 建立 JetStream context
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	// 3. 建立 Stream（若已存在會報錯，可忽略）
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "EVENTS",             // Stream 名稱
		Subjects: []string{"events.*"}, // 要保存的 subject
		Storage:  nats.FileStorage,     // 訊息存檔到磁碟（避免重啟遺失）
		Replicas: 1,                    // 複本數（叢集可設多份）
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		log.Fatal(err)
	}

	// 4. 發布訊息
	for i := 1; i <= 5; i++ {
		ack, err := js.Publish("events.created", []byte(fmt.Sprintf("msg-%d", i)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("[Publish] Message %d stored at sequence %d\n", i, ack.Sequence)
	}

	// 5. Push 訂閱範例
	// 伺服器主動推送訊息
	_, err = js.Subscribe("events.*", func(m *nats.Msg) {
		fmt.Printf("[Push] 收到訊息: %s\n", string(m.Data))
		// 訊息處理成功 → Ack
		m.Ack()
	}, nats.Durable("PUSH_CONSUMER"), nats.ManualAck())
	if err != nil {
		log.Fatal(err)
	}

	// 6. Pull 訂閱範例
	// 由客戶端主動拉取
	sub, err := js.PullSubscribe("events.*", "PULL_CONSUMER")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			// 每次拉 2 筆，最長等待 2 秒
			msgs, err := sub.Fetch(2, nats.MaxWait(2*time.Second))
			if err == nats.ErrTimeout {
				continue
			}
			if err != nil {
				log.Printf("[Pull] Fetch error: %v", err)
				continue
			}

			for _, msg := range msgs {
				fmt.Printf("[Pull] 收到訊息: %s\n", string(msg.Data))
				// 模擬處理
				time.Sleep(500 * time.Millisecond)
				// Ack，否則會被重送
				msg.Ack()
			}
		}
	}()

	// 讓程式跑一段時間觀察
	select {}
}
