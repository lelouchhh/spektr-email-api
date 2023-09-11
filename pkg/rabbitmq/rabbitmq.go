package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// Initialize new RabbitMQ connection
func NewRabbitMQConn() (*amqp.Connection, error) {
	return amqp.Dial("amqp://admin:u9F2DCWg@185.200.241.2:5672/")

}
