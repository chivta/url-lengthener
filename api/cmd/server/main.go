package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/arvlas/url-lengthener/api/internal/config"
	"github.com/arvlas/url-lengthener/api/internal/logger"
	"github.com/arvlas/url-lengthener/api/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("config")
	}

	logger.Init(cfg.LogLevel)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err = server.New(cfg).Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("server")
	}
}
