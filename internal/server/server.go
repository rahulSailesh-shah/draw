package server

import (
	"context"
	"fmt"
	httpSrv "net/http"
	"os/signal"
	"syscall"
	"time"

	"draw/internal/app"
	"draw/internal/transport/http"
	"draw/pkg/auth"

	"github.com/gin-gonic/gin"
)

type Server struct {
	App *app.App
	ctx context.Context
	Engine *gin.Engine
	httpServer *httpSrv.Server
}


func NewServer(ctx context.Context, app *app.App) (*Server, error) {
	engine := gin.Default()

	authKeys, err := auth.LoadKeys(app.Config.Auth.JwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to load auth keys: %w", err)
	}

	// register routes here
	http.RegisterRoutes(engine, authKeys, app )

	srv := &httpSrv.Server{
		Addr:    fmt.Sprintf(":%d", app.Config.Server.Port),
		Handler: engine,
	}

	return &Server{
		App: app,
		ctx: ctx,
		Engine: engine,
		httpServer: srv,
	}, nil
}


func(s *Server) Run() error {
	done := make(chan bool, 1)
	go s.gracefulShutdown(done)

	s.App.Log.Info(s.ctx, "Starting server", "port", s.App.Config.Server.Port)
	if err := s.httpServer.ListenAndServe(); err != nil && err != httpSrv.ErrServerClosed {
		s.App.Log.Error(s.ctx, "Could not start server", "error", err)
		return fmt.Errorf("could not start server: %w", err)
	}

	<-done
	return nil
}


func (s *Server) gracefulShutdown(done chan bool) {
	ctx, stop := signal.NotifyContext(s.ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	s.App.Log.Info(s.ctx, "Shutting down gracefully, press Ctrl+C again to force")
	defer s.App.DB.Close()
	stop()

	timeout := time.Duration(s.App.Config.Server.GracefulShutdownSec) * time.Second
	ctx, cancel := context.WithTimeout(s.ctx, timeout)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.App.Log.Error(s.ctx, "Server forced to shutdown with error", "error", err)
	}

	done <- true
}
