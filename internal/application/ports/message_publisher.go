package ports

import "context"

// MessagePublisher defines the interface for publishing messages to message queues
type MessagePublisher interface {
	// Publish sends a message to the specified queue or exchange
	Publish(ctx context.Context, queueName string, message []byte) error

	// Close closes the connection to the message broker
	Close() error
}
