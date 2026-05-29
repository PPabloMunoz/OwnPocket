package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ppablomunoz/ownpocket/backend/internal/model"
)

func (h *Handler) GetBudgets(c *gin.Context) {
	userID := c.GetUint("user_id")
	budgets, err := h.service.GetBudgets(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, budgets)
}

func (h *Handler) CreateBudget(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		CategoryID uint    `json:"category_id" binding:"required"`
		Period     string  `json:"period" binding:"required"`
		Amount     float64 `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	budget, err := h.service.CreateBudget(userID, req.CategoryID, req.Period, model.NewAmountFromFloat(req.Amount))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, budget)
}
