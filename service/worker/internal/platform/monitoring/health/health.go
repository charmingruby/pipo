package health

import (
	"github.com/charmingruby/pipo/lib/broker/redis"
	"github.com/charmingruby/pipo/lib/logger"
	"github.com/charmingruby/pipo/lib/persistence/postgres"
	"github.com/gin-gonic/gin"
)

type Health struct {
	router *gin.Engine
	logger *logger.Logger
	db     *postgres.Client
	redis  *redis.Client
}

func NewHealth(router *gin.Engine, logger *logger.Logger, db *postgres.Client, redis *redis.Client) *Health {
	return &Health{router: router, logger: logger, db: db, redis: redis}
}

func (h *Health) RegisterProbes() {
	api := h.router.Group("/api")

	api.GET("/health/live", h.makeLivenessProbeEndpoint())
	api.GET("/health/ready", h.makeReadinessProbeEndpoint())
}
