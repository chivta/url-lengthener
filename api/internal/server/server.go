package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/arvlas/url-lengthener/api/internal/config"
	"github.com/arvlas/url-lengthener/api/internal/handler"
	"github.com/arvlas/url-lengthener/api/internal/repository"
	"github.com/arvlas/url-lengthener/api/internal/service"
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

	err = repository.RunMigrations(s.cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("migrations: %w", err)
	}

	redisAddr := strings.TrimPrefix(s.cfg.RedisURL, "redis://")
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	defer func() { _ = rdb.Close() }()
	err = rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}

	urlRepo := repository.NewURLRepo(pool)
	clickRepo := repository.NewClickRepo(pool)

	svc := service.NewURLService(urlRepo, clickRepo)

	router := handler.NewRouter(s.cfg, svc)

	srv := &http.Server{
		Addr:         ":" + s.cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Info().Str("port", s.cfg.Port).Str("env", s.cfg.Environment).Msg("server starting")
	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("listen")
		}
	}()

	<-ctx.Done()
	log.Info().Msg("server shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
