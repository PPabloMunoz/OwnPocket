package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetCategories(c *gin.Context) {
	userID := c.GetUint("user_id")
	categories, err := h.service.GetCategories(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category, err := h.service.CreateCategory(userID, req.Name, req.Type, req.ParentID, req.Color, req.Icon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, category)
}
