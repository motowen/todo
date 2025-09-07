package model

// SendMessageRequest represents the request body for sending a message to SQS
type SendMessageRequest struct {
	QueueName string `json:"queue_name" binding:"required" example:"my-queue"`
	Message   string `json:"message" binding:"required" example:"Hello World"`
}

// SendMessageResponse represents the response for sending a message to SQS
type SendMessageResponse struct {
	Success   bool   `json:"success" example:"true"`
	MessageID string `json:"message_id,omitempty" example:"12345-67890-abcdef"`
	Error     string `json:"error,omitempty" example:""`
}

// SendMessagesRequest represents the request body for sending multiple messages to SQS
type SendMessagesRequest struct {
	QueueName string   `json:"queue_name" binding:"required" example:"my-queue"`
	Messages  []string `json:"messages" binding:"required" example:"[\"message1\", \"message2\"]"`
}

// SendMessagesResponse represents the response for sending multiple messages to SQS
type SendMessagesResponse struct {
	Success        bool     `json:"success" example:"true"`
	SuccessCount   int      `json:"success_count" example:"2"`
	FailedMessages []string `json:"failed_messages,omitempty" example:"[]"`
	Error          string   `json:"error,omitempty" example:""`
}
