package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/kaungmyathan/golang-qwen-mentor-project/internal/config"
	"github.com/kaungmyathan/golang-qwen-mentor-project/internal/server"
	"github.com/kaungmyathan/golang-qwen-mentor-project/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		// Bootstrap logger for the fatal error before zap is ready.
		tmp, _ := zap.NewProduction()
		tmp.Fatal("invalid configuration", zap.Error(err))
	}

	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		tmp, _ := zap.NewProduction()
		tmp.Fatal("failed to build logger", zap.Error(err))
	}
	defer log.Sync() //nolint:errcheck

	log.Info("starting",
		zap.String("env", cfg.AppEnv),
		zap.String("port", cfg.Port),
		zap.String("log_level", cfg.LogLevel),
	)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	srv := server.New(cfg.Port, log)
	if err := srv.Start(ctx); err != nil {
		log.Error("server exited with error", zap.Error(err))
		os.Exit(1)
	}
}

