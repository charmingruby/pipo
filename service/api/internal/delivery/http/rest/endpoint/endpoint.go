package endpoint

import (
	"github.com/charmingruby/pipo/service/api/internal/core/service"
	"github.com/gin-gonic/gin"
)

type Endpoint struct {
	router  *gin.Engine
	service *service.Service
}

func New(router *gin.Engine, service *service.Service) *Endpoint {
	return &Endpoint{router: router, service: service}
}

func (e *Endpoint) Register() {
	api := e.router.Group("/api")

	api.POST("/sentiment/ingest", e.makeIngestRawDataEndpoint())
}
