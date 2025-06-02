// Package health provides the health checks for the application.
package health

import (
	"github.com/charmingruby/pipo-lib/broker/redis"
	"github.com/charmingruby/pipo-lib/logger"
	"github.com/gin-gonic/gin"
)

// Health handles the health checks for the application.
type Health struct {
	router *gin.Engine
	logger *logger.Logger
	redis  *redis.Client
}

// NewHealth constructs a new Health.
//
// router is the router for the application.
// logger is the logger for the application.
// redis is the redis client for the application.
//
// Returns a new Health.
func NewHealth(router *gin.Engine, logger *logger.Logger, redis *redis.Client) *Health {
	return &Health{router: router, logger: logger, redis: redis}
}

// RegisterProbes registers the probes for the application.
//
// Liveness probe is used to check if the application is running.
// Route: /api/health/live.
//
// Readiness probe is used to check if the application is ready to accept requests.
// Route: /api/health/ready.
func (h *Health) RegisterProbes() {
	api := h.router.Group("/api")

	api.GET("/health/live", h.makeLivenessProbeEndpoint())
	api.GET("/health/ready", h.makeReadinessProbeEndpoint())
}
