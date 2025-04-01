package server

import (
	"context"
	"flagon/pkg/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"
)

type HttpServer struct {
	cfg    config.Server
	server *http.Server
	router *gin.Engine
}

type Controller interface {
	Register(router gin.IRouter)
}

func NewHttpServer(cfg config.Server, controllers ...Controller) (*HttpServer, error) {
	err := viper.Sub("http").Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(LogMiddleware, RecoveryMiddleware)
	for _, controller := range controllers {
		controller.Register(router)
	}
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
	return s.server.Shutdown(ctx)
}
