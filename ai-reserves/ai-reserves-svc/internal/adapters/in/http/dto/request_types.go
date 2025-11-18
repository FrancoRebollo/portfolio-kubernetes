package dto

import "time"

type ReqCaptureEvent struct {
	IdEvent      int    `json:"id_event"`
	EventType    string `json:"event_type"`
	EventContent string `json:"event_content"`
}

type ExternalAPIRequest struct {
	Method string            `json:"method"` // "GET" or "POST"
	URL    string            `json:"url"`
	Params map[string]string `json:"params,omitempty"` // For GET
	Body   map[string]any    `json:"body,omitempty"`   // For POST
}

type RequestPushEvent struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`       // semántica de dominio
	RoutingKey string      `json:"routingKey"` // para RabbitMQ
	Origin     string      `json:"origin"`
	Timestamp  time.Time   `json:"timestamp"`
	Payload    interface{} `json:"payload"`
}
