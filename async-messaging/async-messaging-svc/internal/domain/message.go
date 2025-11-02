package domain

import (
	"time"
)

type Event struct {
	ID         string
	Type       string
	RoutingKey string
	Origin     string
	Timestamp  time.Time
	Payload    interface{}
}
