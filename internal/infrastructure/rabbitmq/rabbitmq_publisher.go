package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// RabbitMQPublisher implements the MessagePublisher interface for RabbitMQ
type RabbitMQPublisher struct {
	client  *RabbitMQClient
	channel *amqp091.Channel
}

// NewRabbitMQPublisher creates a new RabbitMQ message publisher
func NewRabbitMQPublisher(client *RabbitMQClient) (*RabbitMQPublisher, error) {
	// Attempt to create a channel, but do not fail if RabbitMQ is down
	channel, _ := client.CreateChannel()

	r := &RabbitMQPublisher{
		client:  client,
		channel: channel, // may be nil
	}

	// Start monitor for channel close and recreation
	go r.monitorChannel()

	return r, nil
}

// Publish sends a message to the specified queue
func (r *RabbitMQPublisher) Publish(ctx context.Context, queueName string, message []byte) error {
	if r.channel == nil || r.channel.IsClosed() {
		// Attempt to get a channel lazily
		ch, err := r.client.CreateChannel()
		if err != nil {
			return fmt.Errorf("publisher channel unavailable: %w", err)
		}
		r.channel = ch
	}

	// Declare the queue (idempotent operation)
	if err := r.client.DeclareQueue(r.channel, queueName); err != nil {
		return err
	}

	cfg := r.client.GetConfig()

	// Publish the message
	err := r.channel.PublishWithContext(
		ctx,
		"",        // exchange (empty string means default exchange)
		queueName, // routing key (queue name)
		false,     // mandatory
		false,     // immediate
		amqp091.Publishing{
			DeliveryMode: getDeliveryMode(cfg.Durable),
			ContentType:  "application/json",
			Body:         message,
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message to queue %s: %w", queueName, err)
	}

	log.Printf("Message published to queue: %s", queueName)
	return nil
}

// getDeliveryMode returns the appropriate delivery mode based on durability
func getDeliveryMode(durable bool) uint8 {
	if durable {
		return amqp091.Persistent // 2
	}
	return amqp091.Transient // 1
}

// monitorChannel watches the publisher channel and re-creates it on close
func (r *RabbitMQPublisher) monitorChannel() {
	for {
		ch := r.channel
		if ch == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		closeCh := ch.NotifyClose(make(chan *amqp091.Error))
		if err := <-closeCh; err != nil {
			log.Printf("Publisher channel closed: %v. Reconnecting...", err)
		} else {
			log.Printf("Publisher channel closed cleanly")
			return
		}

		// Try to recreate channel until success
		for {
			time.Sleep(2 * time.Second)
			newCh, err := r.client.CreateChannel()
			if err != nil {
				log.Printf("Failed to recreate publisher channel: %v", err)
				continue
			}
			r.channel = newCh
			log.Printf("Publisher channel recreated successfully")
			break
		}
	}
}

// Close closes the publisher channel (connection is managed by RabbitMQClient)
func (r *RabbitMQPublisher) Close() error {
	if r.channel != nil && !r.channel.IsClosed() {
		if err := r.channel.Close(); err != nil {
			log.Printf("error closing publisher channel: %v", err)
			return err
		}
	}
	return nil
}
