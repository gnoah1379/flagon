package server

import (
	"context"
	_ "flagon/docs"
	v1 "flagon/pkg/api/v1"
	"flagon/pkg/config"
	"flagon/ui"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type HttpServer struct {
	cfg    config.Server
	server *http.Server
	router *gin.Engine
}

type Controller interface {
	Register(router gin.IRouter)
}

func NewHttpServer(v1Api v1.API) (*HttpServer, error) {
	cfg := config.GetConfig().Server
	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(LogMiddleware, RecoveryMiddleware)
	v1Api.Register(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	ui.Register(router)
	srv := &HttpServer{
		cfg:    cfg,
		server: server,
		router: router,
	}

	return srv, nil
}

func (s *HttpServer) Start() error {
	s.server.Handler = s.router
	slog.Info("Http server starting...", "addr", s.server.Addr)
	if s.cfg.EnableTLS {
		return s.server.ListenAndServeTLS(s.cfg.CertFile, s.cfg.KeyFile)
	}
	return s.server.ListenAndServe()
}

func (s *HttpServer) Stop(ctx context.Context) error {
	slog.Info("Http server shutting down...")
	return s.server.Shutdown(ctx)
}
