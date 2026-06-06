package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/arvlas/url-lengthener/api/internal/config"
	"github.com/arvlas/url-lengthener/api/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err = server.New(cfg).Run(ctx)
	if err != nil {
		log.Fatalf("server: %v", err)
	}
}
