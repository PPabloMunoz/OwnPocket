package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ppablomunoz/ownpocket/backend/internal/model"
)

func (h *Handler) GetTransactions(c *gin.Context) {
	userID := c.GetUint("user_id")
	txs, err := h.service.GetTransactions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, txs)
}

func (h *Handler) CreateTransaction(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		AccountID   uint    `json:"account_id" binding:"required"`
		Amount      float64 `json:"amount" binding:"required"`
		Type        string  `json:"type" binding:"required,oneof=income expense transfer"`
		Date        string  `json:"date" binding:"required"`
		Description string  `json:"description"`
		CategoryID  *uint   `json:"category_id"`
		ToAccountID *uint   `json:"to_account_id"`
		Notes       *string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}
	tx, err := h.service.CreateTransaction(userID, req.AccountID, model.NewAmountFromFloat(req.Amount), req.Type, date, req.Description, req.CategoryID, req.ToAccountID, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, tx)
}

func (h *Handler) GetTransaction(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}
	tx, err := h.service.GetTransaction(userID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}
	c.JSON(http.StatusOK, tx)
}

func (h *Handler) UpdateTransaction(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}
	var updates map[string]any
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx, err := h.service.UpdateTransaction(userID, uint(id), updates)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}
	c.JSON(http.StatusOK, tx)
}

func (h *Handler) DeleteTransaction(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}
	if err := h.service.DeleteTransaction(userID, uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
