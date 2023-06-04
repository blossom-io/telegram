package service

import (
	"blossom/internal/config"
	"blossom/internal/infrastructure/repository"
	"blossom/pkg/logger"
)

type Servicer interface {
	Personer
	Tokener
	Inviter
}

type service struct {
	log  logger.Logger
	cfg  *config.Config
	repo repository.Repository
}

func New(log logger.Logger, cfg *config.Config, repo repository.Repository) Servicer {
	return &service{
		log:  log,
		cfg:  cfg,
		repo: repo,
	}
}
