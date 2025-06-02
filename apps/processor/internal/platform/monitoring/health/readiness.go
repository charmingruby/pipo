package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Health) makeReadinessProbeEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := h.db.Ping(ctx); err != nil {
			h.logger.Error("database is not ready", "error", err)

			c.Status(http.StatusServiceUnavailable)
			return
		}

		if err := h.redis.Ping(ctx); err != nil {
			h.logger.Error("redis is not ready", "error", err)

			c.Status(http.StatusServiceUnavailable)
			return
		}

		c.Status(http.StatusOK)
	}
}
