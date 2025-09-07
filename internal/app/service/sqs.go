package service

import (
	"context"
	"fmt"
	"time"

	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/aws/sqs"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/logger"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/model"
	modelHttp "viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/model/http"
)

// SendMessage sends a single message to the specified SQS queue
func SendMessage(ctx context.Context, request modelHttp.SendMessageRequest) (*modelHttp.SendMessageResponse, model.ServiceResp) {
	logger.Info.Printf("Sending message to queue: %s", request.QueueName)

	// Send message
	err := sqs.TestSQS.SendMessage(ctx, request.Message)
	if err != nil {
		logger.Error.Printf("Failed to send message to queue %s: %v", request.QueueName, err)
		return nil, model.ServiceError.InternalServiceError(fmt.Sprintf("Failed to send message: %v", err))
	}

	response := &modelHttp.SendMessageResponse{
		Success:   true,
		MessageID: generateMessageID(), // You might want to get actual message ID from AWS response
	}

	logger.Info.Printf("Successfully sent message to queue: %s", request.QueueName)
	return response, model.ServiceError.OK
}

// SendMessages sends multiple messages to the specified SQS queue
func SendMessages(ctx context.Context, request modelHttp.SendMessagesRequest) (*modelHttp.SendMessagesResponse, model.ServiceResp) {
	logger.Info.Printf("Sending %d messages to queue: %s", len(request.Messages), request.QueueName)

	if len(request.Messages) == 0 {
		return nil, model.ServiceError.BadRequestError("No messages provided")
	}

	var failedMessages []string
	successCount := 0

	// Send each message individually
	for i, message := range request.Messages {
		err := sqs.TestSQS.SendMessage(ctx, message)
		if err != nil {
			logger.Error.Printf("Failed to send message %d to queue %s: %v", i, request.QueueName, err)
			failedMessages = append(failedMessages, message)
		} else {
			successCount++
		}
	}

	response := &modelHttp.SendMessagesResponse{
		Success:        len(failedMessages) == 0,
		SuccessCount:   successCount,
		FailedMessages: failedMessages,
	}

	if len(failedMessages) > 0 {
		response.Error = fmt.Sprintf("Failed to send %d out of %d messages", len(failedMessages), len(request.Messages))
		logger.Warn.Printf("Partially failed to send messages to queue %s: %d failed, %d succeeded", request.QueueName, len(failedMessages), successCount)
	} else {
		logger.Info.Printf("Successfully sent all %d messages to queue: %s", len(request.Messages), request.QueueName)
	}

	return response, model.ServiceError.OK
}

// generateMessageID generates a simple message ID
// In production, you might want to get the actual message ID from AWS SQS response
func generateMessageID() string {
	// This is a simple implementation - you might want to use UUID or get actual AWS message ID
	return fmt.Sprintf("msg-%d", getCurrentTimestamp())
}

// getCurrentTimestamp returns current unix timestamp
func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}
