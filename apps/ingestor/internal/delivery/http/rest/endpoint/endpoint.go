// Package endpoint provides the endpoints for the application.
package endpoint

import (
	"github.com/charmingruby/pipo/apps/ingestor/internal/core/service"
	"github.com/gin-gonic/gin"
)

// Endpoint manages the registration of the endpoints for the application.
type Endpoint struct {
	router  *gin.Engine
	service *service.Service
}

// New constructs a new Endpoint.
//
// router is the router for the application.
// service is the service for the application.
//
// Returns a new Endpoint.
func New(router *gin.Engine, service *service.Service) *Endpoint {
	return &Endpoint{router: router, service: service}
}

// Register registers the endpoints for the application.
//
// This function registers the endpoints for the application.
func (e *Endpoint) Register() {
	api := e.router.Group("/api")

	api.POST("/sentiment/ingest", e.makeIngestRawDataEndpoint())
}
