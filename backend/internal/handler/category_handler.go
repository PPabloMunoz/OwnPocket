package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ppablomunoz/ownpocket/backend/internal/service"
	"github.com/ppablomunoz/ownpocket/backend/internal/utils"
	"gorm.io/gorm"
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

func (h *Handler) UpdateCategory(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid category id")
		return
	}
	var updates map[string]any
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	category, err := h.service.UpdateCategory(userID, uint(id), updates)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "category not found")
		return
	}
	utils.Success(c, http.StatusOK, category)
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid category id")
		return
	}
	if err := h.service.DeleteCategory(userID, uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "category not found")
		} else if errors.Is(err, service.ErrCategoryHasDependencies) {
			utils.Error(c, http.StatusConflict, "cannot delete: category has child categories or transactions")
		} else {
			utils.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
