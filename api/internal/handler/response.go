package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/arvlas/url-lengthener/api/internal/domain"
)

func writeError(c *gin.Context, err error) {
	c.JSON(statusFor(err), gin.H{"code": errorCode(err)})
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
