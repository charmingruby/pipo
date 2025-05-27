package event

import (
	"github.com/charmingruby/pipo/lib/broker"
	"github.com/charmingruby/pipo/lib/logger"
	"github.com/charmingruby/pipo/service/worker/internal/core/service"
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
