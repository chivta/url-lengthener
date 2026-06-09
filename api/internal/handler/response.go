package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/arvlas/url-lengthener/api/internal/domain"
)

func writeError(c *gin.Context, err error) {
	status := statusFor(err)
	code := errorCode(err)

	ev := log.Warn()
	if status == http.StatusInternalServerError {
		ev = log.Error()
	}
	ev.Err(err).
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Int("status", status).
		Msg(code)

	c.JSON(status, gin.H{"code": code})
}

func statusFor(err error) int {
	switch {
	case errors.Is(err, domain.ErrURLNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrSlugConflict):
		return http.StatusConflict
	case errors.Is(err, domain.ErrURLExpired):
		return http.StatusGone
	case errors.Is(err, domain.ErrInvalidURL):
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

func errorCode(err error) string {
	switch {
	case errors.Is(err, domain.ErrURLNotFound):
		return "url_not_found"
	case errors.Is(err, domain.ErrSlugConflict):
		return "slug_conflict"
	case errors.Is(err, domain.ErrURLExpired):
		return "url_expired"
	case errors.Is(err, domain.ErrInvalidURL):
		return "invalid_url"
	default:
		return "internal_error"
	}
}
