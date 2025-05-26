package event

import (
	"github.com/charmingruby/pipo/internal/sentiment/core/service"
	"github.com/charmingruby/pipo/internal/shared/broker"
)

type Handler struct {
	broker  broker.Broker
	topics  TopicInput
	service *service.Service
}

type TopicInput struct {
	SentimentIngested string
}

func New(broker broker.Broker, topics TopicInput, service *service.Service) *Handler {
	return &Handler{broker: broker, topics: topics, service: service}
}

func (h *Handler) Subscribe() error {
	return h.onSentimentIngested()
}
