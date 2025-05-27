package event

import (
	"github.com/charmingruby/pipo/internal/sentiment/core/service"
	"github.com/charmingruby/pipo/internal/shared/broker"
	"github.com/charmingruby/pipo/pkg/logger"
)

type Handler struct {
	logger  *logger.Logger
	broker  broker.Broker
	topics  TopicInput
	service *service.Service
}

type TopicInput struct {
	SentimentIngested string
}

func New(logger *logger.Logger, broker broker.Broker, topics TopicInput, service *service.Service) *Handler {
	return &Handler{logger: logger, broker: broker, topics: topics, service: service}
}

func (h *Handler) Subscribe() []error {
	return h.onSentimentIngested()
}
