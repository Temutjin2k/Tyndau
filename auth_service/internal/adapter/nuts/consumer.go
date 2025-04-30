package rabbitmq

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
	handler   func(event any)
}

// NewConsumer создает нового Consumer
func NewConsumer(amqpURL, queueName string, handler func(event any)) (*Consumer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &Consumer{
		conn:      conn,
		channel:   ch,
		queueName: queueName,
		handler:   handler,
	}, nil
}

// Start начинает слушать очередь
func (c *Consumer) Start() error {
	msgs, err := c.channel.Consume(
		c.queueName,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var event map[string]any
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Failed to unmarshal event: %v", err)
				continue
			}
			c.handler(event)
		}
	}()

	log.Println("Consumer started listening for messages 📥")
	return nil
}

// Close закрывает соединение и канал
func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
