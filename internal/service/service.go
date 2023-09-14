package service

import (
	"blossom/internal/config"
	"blossom/internal/infrastructure/gpt"
	"blossom/internal/infrastructure/repository"
	"blossom/pkg/logger"
)

type Servicer interface {
	Personer
	Tokener
	Inviter
	Downloader
	AIer
}

type service struct {
	log  logger.Logger
	cfg  *config.Config
	repo repository.Repository
	gpt  gpt.GPTer
}

func New(log logger.Logger, cfg *config.Config, repo repository.Repository, gpt gpt.GPTer) Servicer {
	return &service{
		log:  log,
		cfg:  cfg,
		repo: repo,
		gpt:  gpt,
	}
}
