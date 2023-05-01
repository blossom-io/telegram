package main

import (
	"log"

	"blossom/internal/app"
	"blossom/internal/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
