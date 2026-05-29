package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ppablomunoz/ownpocket/backend/internal/utils"
)

func (h *Handler) GetDashboardSummary(c *gin.Context) {
	userID := c.GetUint("user_id")
	summary, err := h.service.GetDashboardSummary(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, summary)
}
