package rest

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	server http.Server
}

func New(host, port string) (*Server, *gin.Engine) {
	router := gin.Default()

	gin.SetMode(gin.ReleaseMode)

	router.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))

	addr := fmt.Sprintf("%s:%s", host, port)

	return &Server{
		server: http.Server{
			WriteTimeout: 10 * time.Second,
			ReadTimeout:  5 * time.Second,
			IdleTimeout:  120 * time.Second,
			Addr:         addr,
			Handler:      router,
		},
	}, router
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
