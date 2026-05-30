package utils

import "github.com/gin-gonic/gin"

type PaginatedData[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{"data": data, "error": nil})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"data": nil, "error": message})
}
