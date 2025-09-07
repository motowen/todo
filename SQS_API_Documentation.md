# SQS API Documentation

This document describes the AWS SQS integration added to the todo application.

## Features

1. **Send Message API**: Send single message to SQS queue
2. **Send Multiple Messages API**: Send multiple messages to SQS queue
3. **SQS Worker**: Background worker to receive and process messages from SQS queue

## API Endpoints

### Send Single Message

**POST** `/sqs/send-message`

Send a single message to the specified SQS queue.

**Request Body:**
```json
{
    "queue_name": "my-queue",
    "message": "Hello World"
}
```

**Response:**
```json
{
    "success": true,
    "message_id": "msg-1234567890",
    "error": ""
}
```

### Send Multiple Messages

**POST** `/sqs/send-messages`

Send multiple messages to the specified SQS queue.

**Request Body:**
```json
{
    "queue_name": "my-queue",
    "messages": ["message1", "message2", "message3"]
}
```

**Response:**
```json
{
    "success": true,
    "success_count": 3,
    "failed_messages": [],
    "error": ""
}
```

## SQS Worker

The application includes a background SQS worker that:

- Automatically starts when the application starts
- Polls the configured SQS queue for messages
- Processes messages using configurable processors
- Supports multiple concurrent workers
- Handles message deletion after successful processing
- Implements retry logic for failed message processing

### Configuration

The worker is configured via environment variables:

- `AWS_SQS_REGION`: AWS region for SQS (default: us-west-2)
- `AWS_SQS_QUEUE_NAME`: Default queue name for the worker (default: default-queue)

### Message Processing

The default message processor simply logs received messages. You can implement custom processors by:

1. Creating a struct that implements the `MessageProcessor` interface
2. Adding it to the worker service setup in `main.go`

Example custom processor:
```go
type CustomMessageProcessor struct{}

func (p *CustomMessageProcessor) ProcessMessage(ctx context.Context, message sqs.Message) error {
    // Parse message
    var data map[string]interface{}
    if err := json.Unmarshal([]byte(message.Body), &data); err != nil {
        return fmt.Errorf("failed to parse message: %w", err)
    }
    
    // Process the message
    // ... your business logic here ...
    
    return nil
}
```

## Setup Instructions

1. **Configure AWS Credentials**: Set up AWS credentials via AWS CLI, environment variables, or IAM roles
2. **Create SQS Queue**: Create the SQS queue in your AWS account
3. **Set Environment Variables**: Configure the required environment variables (see `env.example`)
4. **Start Application**: The SQS worker will start automatically with the application

## Testing

You can test the SQS functionality by:

1. Starting the application
2. Sending messages via the API endpoints
3. Checking the application logs to see the worker processing messages

Example curl commands:

```bash
# Send single message
curl -X POST http://localhost:8080/sqs/send-message \
  -H "Content-Type: application/json" \
  -d '{"queue_name": "test-queue", "message": "Hello World"}'

# Send multiple messages
curl -X POST http://localhost:8080/sqs/send-messages \
  -H "Content-Type: application/json" \
  -d '{"queue_name": "test-queue", "messages": ["msg1", "msg2", "msg3"]}'
```

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client App    │───▶│   HTTP API      │───▶│   SQS Queue     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                       │
                                                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Logger        │◀───│  SQS Worker     │◀───│   Message       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

The SQS integration provides a robust messaging system for asynchronous processing in your todo application.
