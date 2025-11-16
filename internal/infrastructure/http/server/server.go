package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pull-request-review/config"
	"pull-request-review/internal/infrastructure/adapters/logger"
	"pull-request-review/internal/infrastructure/adapters/router"
)

type Server struct {
	httpServer      *http.Server
	router          router.Router
	shutdownTimeout time.Duration
	logger          logger.Logger
}

func NewServer(router router.Router, cfg config.ServerConfig, log logger.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + cfg.Port,
			Handler:      router,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		router:          router,
		shutdownTimeout: cfg.ShutdownTimeout,
		logger:          log,
	}
}

func (s *Server) Start() {
	go func() {
		s.logger.Info("HTTP server listening", logger.F("address", s.httpServer.Addr))
		if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) && err != nil {
			s.logger.Error(err, "Failed to start HTTP server")
			os.Exit(1)
		}
	}()
}

func (s *Server) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.logger.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		s.logger.Error(err, "Server forced to shutdown")
	}

	s.logger.Info("Server exited")
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}