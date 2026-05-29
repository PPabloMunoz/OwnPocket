package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ppablomunoz/ownpocket/backend/internal/utils"
)

func (h *Handler) GetAccounts(c *gin.Context) {
	userID := c.GetUint("user_id")
	accounts, err := h.service.GetAccounts(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, accounts)
}

func (h *Handler) CreateAccount(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Type        string  `json:"type" binding:"required,oneof=checking savings credit_card cash investment loan"`
		CurrencyID  uint    `json:"currency_id"`
		Description *string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if req.CurrencyID == 0 {
		req.CurrencyID = 1
	}
	account, err := h.service.CreateAccount(userID, req.Name, req.Type, req.CurrencyID, req.Description)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusCreated, account)
}

func (h *Handler) GetAccount(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid account id")
		return
	}
	account, err := h.service.GetAccount(userID, uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, "account not found")
		return
	}
	utils.Success(c, http.StatusOK, account)
}

func (h *Handler) UpdateAccount(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid account id")
		return
	}
	var updates map[string]any
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	account, err := h.service.UpdateAccount(userID, uint(id), updates)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "account not found")
		return
	}
	utils.Success(c, http.StatusOK, account)
}

func (h *Handler) DeleteAccount(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid account id")
		return
	}
	if err := h.service.DeleteAccount(userID, uint(id)); err != nil {
		utils.Error(c, http.StatusNotFound, "account not found")
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
