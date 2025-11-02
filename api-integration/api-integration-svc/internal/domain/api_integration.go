package domain

import "time"

type ExternalAPIRequest struct {
	Method string
	URL    string
	Params map[string]string
	Body   map[string]any
}

type ExternalAPIResponse struct {
	Status     string
	StatusCode int
	Data       any
}

type Event struct {
	ID         string
	Type       string
	RoutingKey string
	Origin     string
	Timestamp  time.Time
	Payload    interface{}
}
