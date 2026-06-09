package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func loggingMiddleware() gin.HandlerFunc {
	skip := map[string]bool{"/health": true, "/metrics": true}
	return func(c *gin.Context) {
		if skip[c.Request.URL.Path] {
			c.Next()
			return
		}
		start := time.Now()
		c.Next()
		log.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("duration", time.Since(start)).
			Str("request_id", c.GetHeader("X-Request-Id")).
			Msg("request")
	}
}

func corsMiddleware(allowedOrigins string) gin.HandlerFunc {
	set := make(map[string]bool)
	for o := range strings.SplitSeq(allowedOrigins, ",") {
		set[strings.TrimSpace(o)] = true
	}
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if set[origin] || set["*"] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type")
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
