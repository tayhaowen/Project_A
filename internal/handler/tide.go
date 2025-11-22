package handler

import (
	"net/http"

	"project-a/internal/service"

	"github.com/gin-gonic/gin"
)

// TideHandler exposes endpoints for tide timings queries.
type TideHandler struct {
	tideService *service.TideService
}

// NewTideHandler constructs a TideHandler instance.
func NewTideHandler(tideService *service.TideService) *TideHandler {
	return &TideHandler{tideService: tideService}
}

// RegisterRoutes binds the tide routes onto the provided router group/engine.
func (h *TideHandler) RegisterRoutes(router gin.IRoutes) {
	router.GET("/tide-timings", h.getTideTimings)
}

func (h *TideHandler) getTideTimings(c *gin.Context) {
	data, err := h.tideService.GetTideTimings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}
