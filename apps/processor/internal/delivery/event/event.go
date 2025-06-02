// Package event provides the event handlers for the worker service.
package event

import (
	"github.com/charmingruby/pipo-lib/broker"
	"github.com/charmingruby/pipo-lib/logger"
	"github.com/charmingruby/pipo/apps/processor/internal/core/service"
)

// Handler is the handler for the event service.
type Handler struct {
	logger  *logger.Logger
	service *service.Service
	broker  broker.Broker
	topics  TopicInput
}

// TopicInput is the input to handle the topics.
type TopicInput struct {
	// SentimentIngested is the topic of the sentiment ingested.
	SentimentIngested string
}

// New constructs a new Handler.
//
// logger is the logger for the application.
// broker is the broker for the application.
// topics is the topics for the application.
// service is the service for the application.
//
// Returns a new Handler.
func New(logger *logger.Logger, broker broker.Broker, topics TopicInput, service *service.Service) *Handler {
	return &Handler{logger: logger, broker: broker, topics: topics, service: service}
}

// Subscribe subscribes to the topics.
//
// Returns a list of errors.
func (h *Handler) Subscribe() []error {
	return h.onSentimentIngested()
}
