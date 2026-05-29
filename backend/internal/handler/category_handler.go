package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ppablomunoz/ownpocket/backend/internal/utils"
)

func (h *Handler) GetCategories(c *gin.Context) {
	userID := c.GetUint("user_id")
	categories, err := h.service.GetCategories(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusOK, categories)
}

func (h *Handler) CreateCategory(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		Name     string  `json:"name" binding:"required"`
		Type     string  `json:"type" binding:"required,oneof=income expense"`
		ParentID *uint   `json:"parent_id"`
		Color    *string `json:"color"`
		Icon     *string `json:"icon"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	category, err := h.service.CreateCategory(userID, req.Name, req.Type, req.ParentID, req.Color, req.Icon)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.Success(c, http.StatusCreated, category)
}
