package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/aws/sqs"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/config"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/logger"
)

// MessageProcessor defines the interface for processing SQS messages
type MessageProcessor interface {
	ProcessMessage(ctx context.Context, message sqs.Message) error
}

// SQSWorker represents a worker that processes SQS messages
type SQSWorker struct {
	queueName    string
	sqsManager   sqs.SQSAPI
	processor    MessageProcessor
	pollInterval time.Duration
	maxRetries   int
	workerCount  int
	running      bool
	stopChan     chan struct{}
	wg           sync.WaitGroup
	mu           sync.RWMutex
}

// SQSWorkerConfig holds configuration for SQS worker
type SQSWorkerConfig struct {
	QueueName    string
	Processor    MessageProcessor
	PollInterval time.Duration // How long to wait between polls when no messages
	MaxRetries   int           // Maximum retries for failed message processing
	WorkerCount  int           // Number of concurrent workers
}

// NewSQSWorker creates a new SQS worker
func NewSQSWorker(cfg SQSWorkerConfig) (*SQSWorker, error) {
	if cfg.QueueName == "" {
		return nil, fmt.Errorf("queue name is required")
	}
	if cfg.Processor == nil {
		return nil, fmt.Errorf("message processor is required")
	}
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = 5 * time.Second
	}
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 3
	}
	if cfg.WorkerCount <= 0 {
		cfg.WorkerCount = 1
	}

	sqsManager, err := sqs.NewBaseManager(sqs.Config{
		QueueName: cfg.QueueName,
		Region:    config.Env.AWSSQSRegion,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create SQS manager: %w", err)
	}

	return &SQSWorker{
		queueName:    cfg.QueueName,
		sqsManager:   &sqsManager,
		processor:    cfg.Processor,
		pollInterval: cfg.PollInterval,
		maxRetries:   cfg.MaxRetries,
		workerCount:  cfg.WorkerCount,
		stopChan:     make(chan struct{}),
	}, nil
}

// Start starts the SQS worker
func (w *SQSWorker) Start(ctx context.Context) {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		logger.Warn.Printf("SQS worker for queue %s is already running", w.queueName)
		return
	}
	w.running = true
	w.mu.Unlock()

	logger.Info.Printf("Starting SQS worker for queue %s with %d workers", w.queueName, w.workerCount)

	// Start multiple worker goroutines
	for i := 0; i < w.workerCount; i++ {
		w.wg.Add(1)
		go w.workerLoop(ctx, i)
	}
}

// Stop stops the SQS worker gracefully
func (w *SQSWorker) Stop() {
	w.mu.Lock()
	if !w.running {
		w.mu.Unlock()
		return
	}
	w.running = false
	w.mu.Unlock()

	logger.Info.Printf("Stopping SQS worker for queue %s", w.queueName)

	// Signal all workers to stop
	close(w.stopChan)

	// Wait for all workers to finish
	w.wg.Wait()

	logger.Info.Printf("SQS worker for queue %s stopped", w.queueName)
}

// IsRunning returns true if the worker is currently running
func (w *SQSWorker) IsRunning() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.running
}

// workerLoop is the main loop for each worker goroutine
func (w *SQSWorker) workerLoop(ctx context.Context, workerID int) {
	defer w.wg.Done()

	logger.Info.Printf("SQS worker %d started for queue %s", workerID, w.queueName)

	for {
		select {
		case <-w.stopChan:
			logger.Info.Printf("SQS worker %d stopping for queue %s", workerID, w.queueName)
			return
		case <-ctx.Done():
			logger.Info.Printf("SQS worker %d context cancelled for queue %s", workerID, w.queueName)
			return
		default:
			w.processMessages(ctx, workerID)
		}
	}
}

// processMessages polls for and processes messages
func (w *SQSWorker) processMessages(ctx context.Context, workerID int) {
	// Poll for messages with long polling (20 seconds max)
	hasMessage, message, err := w.sqsManager.ReceiveMessage(ctx, 20, 300) // 20s wait, 5min visibility timeout
	if err != nil {
		logger.Error.Printf("SQS worker %d failed to receive message from queue %s: %v", workerID, w.queueName, err)
		time.Sleep(w.pollInterval)
		return
	}

	if !hasMessage {
		// No messages available, short sleep before next poll
		time.Sleep(1 * time.Second)
		return
	}

	logger.Info.Printf("SQS worker %d received message from queue %s", workerID, w.queueName)

	// Process the message with retries
	if w.processMessageWithRetries(ctx, message, workerID) {
		// Successfully processed, delete the message
		err = w.sqsManager.DeleteMessage(ctx, message.ReceiptHandle)
		if err != nil {
			logger.Error.Printf("SQS worker %d failed to delete message from queue %s: %v", workerID, w.queueName, err)
		} else {
			logger.Info.Printf("SQS worker %d successfully processed and deleted message from queue %s", workerID, w.queueName)
		}
	} else {
		logger.Error.Printf("SQS worker %d failed to process message from queue %s after %d retries", workerID, w.queueName, w.maxRetries)
		// Message will become visible again after visibility timeout expires
	}
}

// processMessageWithRetries processes a message with retry logic
func (w *SQSWorker) processMessageWithRetries(ctx context.Context, message sqs.Message, workerID int) bool {
	for attempt := 1; attempt <= w.maxRetries; attempt++ {
		err := w.processor.ProcessMessage(ctx, message)
		if err == nil {
			return true
		}

		logger.Warn.Printf("SQS worker %d failed to process message (attempt %d/%d) from queue %s: %v",
			workerID, attempt, w.maxRetries, w.queueName, err)

		if attempt < w.maxRetries {
			// Wait before retry with exponential backoff
			backoffDuration := time.Duration(attempt) * time.Second
			time.Sleep(backoffDuration)
		}
	}
	return false
}

// DefaultMessageProcessor is a simple implementation of MessageProcessor for demonstration
type DefaultMessageProcessor struct{}

// ProcessMessage processes a message (default implementation just logs it)
func (p *DefaultMessageProcessor) ProcessMessage(ctx context.Context, message sqs.Message) error {
	logger.Info.Printf("Processing message: %s", message.Body)

	// Add your actual message processing logic here
	// For example: parse JSON, call business logic, update database, etc.

	return nil
}
