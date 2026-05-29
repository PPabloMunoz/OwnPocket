package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetDashboardSummary(c *gin.Context) {
	userID := c.GetUint("user_id")
	summary, err := h.service.GetDashboardSummary(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}
