package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Message struct {
	Body          string
	ReceiptHandle string
}

type SQSAPI interface {
	SendMessage(ctx context.Context, body string) error
	ReceiveMessage(ctx context.Context, waitTime, visibilityTimeout int32) (hasMessage bool, message Message, err error)
	DeleteMessage(ctx context.Context, receiptHandle string) error
}

var (
	TestSQS SQSAPI
)

const MaxVisibilityTimeout int32 = 43200 - 5

type BaseSQSAPI struct {
	queueURL *string
	client   *sqs.Client
}

type Config struct {
	QueueName string
	Region    string
}

func NewBaseManager(setupConfig Config) (BaseSQSAPI, error) {
	ctx := context.Background()
	manager := BaseSQSAPI{}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return manager, err
	}

	client := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		o.Region = setupConfig.Region
	})

	result, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: &setupConfig.QueueName,
	})
	if err != nil {
		return manager, err
	}

	manager = BaseSQSAPI{
		queueURL: result.QueueUrl,
		client:   client,
	}
	return manager, nil
}

func (manager *BaseSQSAPI) SendMessage(ctx context.Context, body string) error {
	_, err := manager.client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(body),
		QueueUrl:    manager.queueURL,
	})
	return err
}

func (manager *BaseSQSAPI) ReceiveMessage(ctx context.Context, waitTimeSeconds, visibilityTimeout int32) (bool, Message, error) {
	msgOutput, err := manager.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl:            manager.queueURL,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     waitTimeSeconds,
		VisibilityTimeout:   min(visibilityTimeout, MaxVisibilityTimeout),
	})
	if err != nil {
		return false, Message{}, err
	}
	if len(msgOutput.Messages) == 0 {
		return false, Message{}, nil
	} else {
		return true, Message{Body: *msgOutput.Messages[0].Body, ReceiptHandle: *msgOutput.Messages[0].ReceiptHandle}, nil
	}
}

func (manager *BaseSQSAPI) DeleteMessage(ctx context.Context, receiptHandle string) error {
	_, err := manager.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      manager.queueURL,
		ReceiptHandle: aws.String(receiptHandle),
	})
	return err
}
