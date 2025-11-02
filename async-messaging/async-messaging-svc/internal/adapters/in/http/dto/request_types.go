package dto

import (
	"time"
)

type RequestPushEvent struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`       // semántica de dominio
	RoutingKey string      `json:"routingKey"` // para RabbitMQ
	Origin     string      `json:"origin"`
	Timestamp  time.Time   `json:"timestamp"`
	Payload    interface{} `json:"payload"`
}
