package service

import (
	"github.com/charmingruby/pipo/internal/shared/messaging"
	"github.com/charmingruby/pipo/pkg/logger"
)

type Service struct {
	logger *logger.Logger
	broker messaging.Broker
}

func NewService(logger *logger.Logger, broker messaging.Broker) *Service {
	return &Service{
		logger: logger,
		broker: broker,
	}
}
