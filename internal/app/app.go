package app

import (
	"context"
	"os"
	"os/signal"

	"blossom/internal/config"
	"blossom/internal/infrastructure/bot"
	"blossom/internal/infrastructure/repository"
	"blossom/internal/service"
	"blossom/pkg/logger"
	"blossom/pkg/postgres"
)

// Run runs application.
func Run(cfg *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log := logger.New(cfg.Bot.LogLevel)
	log.Info("starting bot...")

	// Database
	DB, err := postgres.New(cfg.Connections.Postgres.URL)
	if err != nil {
		log.Fatal("app - Run - postgres.New - error initializing database: %w", err)
	}
	defer DB.Close()

	// Repositories
	repo := repository.New(DB)

	svc := service.New(log, cfg, repo)

	botSvc, err := bot.New(ctx, cfg, log, svc)
	if err != nil {
		log.Error("app - Run - bot.New: %w", err)
		os.Exit(1)
	}

	botSvc.Run(ctx)

	<-ctx.Done()
	log.Info("Gracefully shutting down...")
}
