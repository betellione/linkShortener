package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// PostHandler создает сокращенный URL
// @Summary Создает сокращенный URL
// @Description Создает сокращенный URL из длинного
// @Tags URL
// @Accept  json
// @Produce  json
// @Param   url      body    JsonRequest     true  "URL Request"
// @Success 200 {object} JsonRequest
// @Failure 400 {object} JsonRequest
// @Failure 500 {object} JsonRequest
// @Router /create [post]
func PostHandler(db Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var jsonReq JsonRequest
		if err := c.BindJSON(&jsonReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		shortened, err := db.FindOrCreateURL(jsonReq.URL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or find URL"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"shortId": shortened})
	}
}

// RedirectHandler перенаправляет на оригинальный URL
// @Summary Перенаправляет на оригинальный URL
// @Description Получает короткий URL и перенаправляет на соответствующий полный URL
// @Tags URL
// @Accept  json
// @Produce  json
// @Param   shortURL     path    string     true  "Короткий URL"
// @Success 302 {string} string "Redirected"
// @Failure 400 {object} JsonRequest
// @Failure 404 {object} JsonRequest
// @Failure 500 {object} JsonRequest
// @Router /{shortURL} [get]
func RedirectHandler(db Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortURL := c.Param("shortURL")
		url, err := db.GetURL(shortURL)
		if err != nil {
			c.JSON(getStatusError(err), gin.H{"error": "URL not found"})
			return
		}
		c.Redirect(http.StatusFound, url)
	}
}
