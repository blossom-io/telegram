package app

import (
	"context"
	"os"
	"os/signal"

	"blossom/internal/config"
	"blossom/internal/infrastructure/bot"
	"blossom/internal/infrastructure/gpt"
	"blossom/internal/infrastructure/repository"
	"blossom/internal/service"
	"blossom/pkg/logger"
	"blossom/pkg/postgres"
)

// Run runs application.
func Run(cfg *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log := logger.New()
	log.Info("starting bot...")

	// Database
	DB, err := postgres.New(cfg.Connections.Postgres.URL)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	defer DB.Close()

	// Repositories
	repo := repository.New(DB)

	// AI
	gpt := gpt.New(cfg, log)

	svc := service.New(log, cfg, repo, gpt)

	botSvc, err := bot.New(ctx, cfg, log, svc)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	botSvc.Run(ctx)

	<-ctx.Done()
	log.Info("Gracefully shutting down...")
}
