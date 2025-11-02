package http

import (
	"github.com/FrancoRebollo/async-messaging-svc/internal/ports"
)

type MessageHandler struct {
	serv ports.MessageService
}

func NewMessageHandler(serv ports.MessageService) *MessageHandler {
	return &MessageHandler{
		serv,
	}
}
