package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/arvlas/url-lengthener/api/internal/domain"
)

type urlHandler struct {
	svc domain.URLService
}

func (h *urlHandler) ShortenURL(c *gin.Context) {
	var req struct {
		URL        string `json:"url"         binding:"required"`
		CustomSlug string `json:"custom_slug"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		writeError(c, domain.ErrInvalidURL)
		return
	}
	u, err := h.svc.Shorten(c.Request.Context(), req.URL, req.CustomSlug)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, u)
}

func (h *urlHandler) GetURL(c *gin.Context) {
	u, err := h.svc.Get(c.Request.Context(), c.Param("slug"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *urlHandler) DeleteURL(c *gin.Context) {
	err := h.svc.Delete(c.Request.Context(), c.Param("slug"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *urlHandler) RedirectURL(c *gin.Context) {
	u, err := h.svc.Resolve(c.Request.Context(), c.Param("slug"))
	if err != nil {
		writeError(c, err)
		return
	}
	c.Redirect(http.StatusMovedPermanently, u.OriginalURL)
}
