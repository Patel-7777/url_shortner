package handlers

import (
	"net/http"

	"github.com/drashti/url_shortner/internal/service"
	"github.com/gin-gonic/gin"
)

type ShortenerHandler struct {
	service *service.ShortenerService
}

func NewShortenerHandler(service *service.ShortenerService) *ShortenerHandler {
	return &ShortenerHandler{service: service}
}

type CreateShortURLRequest struct {
	URL string `json:"url" binding:"required"`
}

type CreateShortURLResponse struct {
	ShortCode string `json:"short_code"`
}

func (h *ShortenerHandler) CreateShortURL(c *gin.Context) {
	var req CreateShortURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	shortCode, err := h.service.CreateShortURL(c.Request.Context(), req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL"})
		return
	}

	c.JSON(http.StatusOK, CreateShortURLResponse{
		ShortCode: shortCode,
	})
}

func (h *ShortenerHandler) RedirectToURL(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Short code is required"})
		return
	}

	originalURL, err := h.service.GetOriginalURL(c.Request.Context(), shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get original URL"})
		return
	}

	if originalURL == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}
