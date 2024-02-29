package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"shorten/pkg/service"
	"shorten/pkg/structs"
)

// PostHandler Определение обработчика HTTP POST запросов для создания коротких URL.
func PostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var jsonReq structs.JsonRequest
		if err := c.BindJSON(&jsonReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		shortened, err := service.FindOrCreateURL(db, jsonReq.URL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or find URL"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"shortId": shortened})
	}
}
