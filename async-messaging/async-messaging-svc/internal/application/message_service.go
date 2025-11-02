package application

import (
	"github.com/FrancoRebollo/async-messaging-svc/internal/platform/config"

	"github.com/FrancoRebollo/async-messaging-svc/internal/ports"
)

type MessageService struct {
	hr   ports.MessageRepository
	rmq  ports.MessageQueue
	conf config.App
}

func NewMessageService(hr ports.MessageRepository, rmq ports.MessageQueue, conf config.App) *MessageService {
	return &MessageService{
		hr,
		rmq,
		conf,
	}
}
