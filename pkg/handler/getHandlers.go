package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"shorten/pkg/structs"
)

// RedirectHandler Определение обработчика HTTP GET запросов для перенаправления по короткому URL.
func RedirectHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortURL := c.Param("shortURL")
		var urlRecord structs.ShortURL
		if err := db.Where("short_url = ?", shortURL).First(&urlRecord).Error; err != nil {
			c.JSON(getStatusError(err), gin.H{"error": "URL not found"})
			return
		}
		c.Redirect(http.StatusFound, urlRecord.URL)
	}
}

// GetStatusError Вспомогательная функция для определения статуса ошибки при общении с базой данных.
func getStatusError(err error) (status int) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}
