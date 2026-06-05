package handler

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
		slog.Info("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", time.Since(start),
			"request_id", c.GetHeader("X-Request-Id"),
		)
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
