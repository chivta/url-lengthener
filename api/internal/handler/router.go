package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/arvlas/url-lengthener/api/internal/config"
	"github.com/arvlas/url-lengthener/api/internal/domain"
)

func NewRouter(cfg *config.Config, svc domain.URLService) http.Handler {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(loggingMiddleware())
	r.Use(corsMiddleware(cfg.AllowedOrigins))

	r.GET("/health", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	h := &urlHandler{svc: svc}
	v1 := r.Group("/api/v1")
	v1.POST("/urls/suggest", h.SuggestSlugs)
	v1.POST("/urls", h.ShortenURL)
	v1.GET("/urls/:slug", h.GetURL)
	v1.DELETE("/urls/:slug", h.DeleteURL)

	r.GET("/:slug", h.RedirectURL)

	return r
}
