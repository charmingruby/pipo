package health

import (
	"github.com/charmingruby/pipo/lib/broker/redis"
	"github.com/charmingruby/pipo/lib/logger"
	"github.com/gin-gonic/gin"
)

type Health struct {
	router *gin.Engine
	logger *logger.Logger
	redis  *redis.Client
}

func NewHealth(router *gin.Engine, logger *logger.Logger, redis *redis.Client) *Health {
	return &Health{router: router, logger: logger, redis: redis}
}

func (h *Health) RegisterProbes() {
	api := h.router.Group("/api")

	api.GET("/health/live", h.makeLivenessProbeEndpoint())
	api.GET("/health/ready", h.makeReadinessProbeEndpoint())
}
