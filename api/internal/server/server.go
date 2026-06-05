package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/redis/go-redis/v9"

	"github.com/arvlas/url-shortener/api/internal/config"
	"github.com/arvlas/url-shortener/api/internal/handler"
	"github.com/arvlas/url-shortener/api/internal/repository"
	"github.com/arvlas/url-shortener/api/internal/service"
)

type Server struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Server {
	return &Server{cfg: cfg}
}

func (s *Server) Run(ctx context.Context) error {
	pool, err := repository.NewPool(ctx, s.cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("db pool: %w", err)
	}
	defer pool.Close()

	if err := repository.RunMigrations(pool); err != nil {
		return fmt.Errorf("migrations: %w", err)
	}

	redisAddr := strings.TrimPrefix(s.cfg.RedisURL, "redis://")
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	defer rdb.Close()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}

	urlRepo := repository.NewURLRepo(pool)
	clickRepo := repository.NewClickRepo(pool)

	claudeClient := anthropic.NewClient(option.WithAPIKey(s.cfg.AnthropicAPIKey))
	svc := service.NewURLService(urlRepo, clickRepo, &claudeClient)

	router := handler.NewRouter(s.cfg, svc)

	srv := &http.Server{
		Addr:         ":" + s.cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	slog.Info("server starting", "port", s.cfg.Port, "env", s.cfg.Environment)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("server shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
