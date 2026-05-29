package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ppablomunoz/ownpocket/backend/internal/model"
	"github.com/ppablomunoz/ownpocket/backend/internal/utils"
)

func (h *Handler) GetTransactions(c *gin.Context) {
	userID := c.GetUint("user_id")
	txs, err := h.service.GetTransactions(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, txs)
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
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if req.Type == "transfer" {
		if req.ToAccountID == nil {
			utils.Error(c, http.StatusBadRequest, "to_account_id is required for transfers")
			return
		}
		if *req.ToAccountID == req.AccountID {
			utils.Error(c, http.StatusBadRequest, "source and destination accounts must be different")
			return
		}
	}
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid date format, use YYYY-MM-DD")
		return
	}
	tx, err := h.service.CreateTransaction(userID, req.AccountID, model.NewAmountFromFloat(req.Amount), req.Type, date, req.Description, req.CategoryID, req.ToAccountID, req.Notes)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusCreated, tx)
}

func (h *Handler) GetTransaction(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid transaction id")
		return
	}
	tx, err := h.service.GetTransaction(userID, uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, "transaction not found")
		return
	}
	utils.Success(c, http.StatusOK, tx)
}

func (h *Handler) UpdateTransaction(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid transaction id")
		return
	}
	var updates map[string]any
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if txType, ok := updates["type"].(string); ok && txType == "transfer" {
		if toID, ok := updates["to_account_id"]; !ok || toID == nil {
			utils.Error(c, http.StatusBadRequest, "to_account_id is required for transfers")
			return
		}
		if accID, ok := updates["account_id"]; ok {
			toID := uint(updates["to_account_id"].(float64))
			fromID := uint(accID.(float64))
			if toID == fromID {
				utils.Error(c, http.StatusBadRequest, "source and destination accounts must be different")
				return
			}
		}
	}
	tx, err := h.service.UpdateTransaction(userID, uint(id), updates)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "transaction not found")
		return
	}
	utils.Success(c, http.StatusOK, tx)
}

func (h *Handler) DeleteTransaction(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid transaction id")
		return
	}
	if err := h.service.DeleteTransaction(userID, uint(id)); err != nil {
		utils.Error(c, http.StatusNotFound, "transaction not found")
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
