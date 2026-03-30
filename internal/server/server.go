package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

const shutdownTimeout = 30 * time.Second

// Server wraps the HTTP server and its dependencies.
type Server struct {
	http   *http.Server
	logger *zap.Logger
}

// New creates a Server with all routes registered.
func New(port string, logger *zap.Logger) *Server {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(zapMiddleware(logger))
	r.Use(middleware.Recoverer)

	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		handleReady(w, r.WithContext(r.Context()))
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		handleHealth(w, r.WithContext(r.Context()))
	})

	return &Server{
		http: &http.Server{
			Addr:         ":" + port,
			Handler:      r,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		logger: logger,
	}
}

// Start begins listening and blocks until the context is cancelled, then drains
// in-flight requests within shutdownTimeout.
func (s *Server) Start(ctx context.Context) error {
	listenErr := make(chan error, 1)

	go func() {
		s.logger.Info("server listening", zap.String("addr", s.http.Addr))
		if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			listenErr <- fmt.Errorf("listen: %w", err)
		}
		close(listenErr)
	}()

	select {
	case err := <-listenErr:
		return err
	case <-ctx.Done():
		s.logger.Info("shutdown signal received, draining connections")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := s.http.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown: %w", err)
	}

	s.logger.Info("server stopped cleanly")
	return nil
}

// zapMiddleware logs every request using the provided zap logger.
func zapMiddleware(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()
			next.ServeHTTP(ww, r)
			log.Info("request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.Status()),
				zap.Duration("latency", time.Since(start)),
				zap.String("request_id", middleware.GetReqID(r.Context())),
			)
		})
	}
}
