package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go-base/internal/app/router"
	"go-base/internal/pkg/aws/s3"
	"go-base/internal/pkg/aws/sqs"
	"go-base/internal/pkg/config"
	"go-base/internal/pkg/database"
	"go-base/internal/pkg/http/client"
	"go-base/internal/pkg/logger"
)

func Setup() {
	var err error

	if err = config.Setup(); err != nil {
		log.Fatal(err)
	}

	/*
		if err = cache.GetInstance().Setup(cache.Config{
			Type:         config.Env.RedisType,
			EndpointList: config.Env.RedisEndpointList,
			Password:     config.Env.RedisPassword,
		}); err != nil {
			log.Fatalf("cache Setup, error:%v", err)
		}
	*/

	if err = database.Setup(config.Env.MongoURI); err != nil {
		log.Fatalf("database Setup, error:%v", err)
	}

	/*
		if err = postgres.GetInstance().Setup(postgres.Config{
			Username:                config.Env.PostgresUsername,
			Password:                config.Env.PostgresPassword,
			Host:                    config.Env.PostgresHost,
			Port:                    config.Env.PostgresPort,
			TableName:               config.Env.PostgresName,
			MinConnSize:             config.Env.PostgresMinConnSize,
			MaxConnSize:             config.Env.PostgresMaxConnSize,
			MaxConnIdleTimeBySecond: time.Duration(config.Env.PostgresMaxConnIdleTimeBySecond),
			MaxConnLifetimeBySecond: time.Duration(config.Env.PostgresMaxConnLifeTimeBySecond),
		}); err != nil {
			log.Fatalf("postgres Setup, error:%v", err)
		}
	*/

	/*
		if err = queue.GetInstance().Setup(queue.Config{
			Url: config.Env.NatsUrl,
		}); err != nil {
			log.Fatalf("queue Setup, error:%v", err)
		}
	*/

	/*
		if err = search.GetInstance().Setup(search.Config{
			Url:         config.Env.ElasticsearchUrl,
			IndexPrefix: config.Env.ElasticsearchIndexPrefix,
		}); err != nil {
			log.Fatalf("search Setup, error:%v", err)
		}
	*/

	if err = logger.Setup(config.Env.LogLevel); err != nil {
		log.Fatal(err)
	}

	client.Setup()

	if s3API, err := s3.NewBaseS3API(s3.Config{
		AWSS3Bucket:         config.Env.AWSS3Bucket,
		AWSS3Region:         config.Env.AWSS3Region,
		IsEnabledAccelerate: config.Env.IsEnabledAccelerate,
	}); err != nil {
		log.Fatalf("aws Setup, error:%v", err)
	} else {
		s3.SetInstance(s3API)
	}

	if sqsTest, err := sqs.NewBaseManager(sqs.Config{
		QueueName: config.Env.AWSSQSQueueName,
		Region:    config.Env.AWSSQSRegion,
	}); err != nil {
		log.Fatalf("sqs LodCreated Setup, region: %s, queue name: %s, error:%v", config.Env.AWSSQSRegion, config.Env.AWSSQSQueueName, err)
	} else {
		sqs.TestSQS = &sqsTest
	}

	if err = router.Setup(); err != nil {
		log.Fatal(err)
	}
}

func Close() {
}

func RunServer() {
	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Env.Port),
		Handler:      router.Router,
		ReadTimeout:  30 * time.Minute,
		WriteTimeout: 30 * time.Minute,
	}
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("%s\n", err)
	}
}

func receiveMessage() {
	for {
		hasMsg, msg, err := sqs.TestSQS.ReceiveMessage(context.TODO(), 20, 30)
		if err != nil {
			log.Fatalf("ReceiveMessage, error:%v", err)
		} else if hasMsg {
			go func() {
				logger.Info.Printf("received sqs message: %v", msg)

				defer func() {
					logger.Info.Printf("delete sqs message: %v", msg)
					if err := sqs.TestSQS.DeleteMessage(context.TODO(), msg.ReceiptHandle); err != nil {
						err = fmt.Errorf("fetch lod file failed to delete sqs message: %v", err)
						logger.Error.Printf(err.Error())
					}
				}()
			}()
		} else {
			logger.Info.Printf("no message received")
		}
	}
}

// @title        Community Service Swagger
// @description  this service is Community Service
func main() {
	Setup()
	defer Close()
	go receiveMessage()
	RunServer()
}
