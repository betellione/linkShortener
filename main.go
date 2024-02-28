package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type ShortURL struct {
	gorm.Model
	ShortURL string `gorm:"unique_index"`
	URL      string `gorm:"unique_index"`
}

type JsonRequest struct {
	URL string `json:"url"`
}

func To62Base(num uint) string {
	digits := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	if num == 0 {
		return "0"
	}

	base := uint(62)
	result := ""
	for num > 0 {
		remainder := num % base
		num /= base
		result = string(digits[remainder]) + result
	}

	return result
}

func (base *ShortURL) AfterCreate(tx *gorm.DB) (err error) {
	base.ShortURL = To62Base(base.ID)
	return tx.Model(base).Update("ShortURL", base.ShortURL).Error
}

func findOrCreateURL(db *gorm.DB, inputURL string) (string, error) {
	var urlRecord ShortURL
	urlRecord.URL = inputURL
	err := db.Where("url = ?", inputURL).FirstOrCreate(&urlRecord).Error
	return urlRecord.ShortURL, err
}

func getDatabaseConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=localhost user=%s dbname=%s password=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("DB_PORT"))
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func shortenURLHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var jsonReq JsonRequest
		if err := c.BindJSON(&jsonReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		shortened, err := findOrCreateURL(db, jsonReq.URL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"shortURL": shortened})
	}
}

func redirectHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortURL := c.Param("shortURL")
		var urlRecord ShortURL
		if err := db.Where("short_url = ?", shortURL).First(&urlRecord).Error; err != nil {
			status := http.StatusInternalServerError
			message := "Error retrieving the original URL"
			if errors.Is(err, gorm.ErrRecordNotFound) {
				status = http.StatusNotFound
				message = "URL not found"
			}
			c.JSON(status, gin.H{"error": message})
			return
		}
		c.Redirect(http.StatusFound, urlRecord.URL)
	}
}

func getGinRouter() (*gin.Engine, error) {
	gin.SetMode(os.Getenv("GIN_MODE"))
	router := gin.Default()
	err := router.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		return nil, err
	}
	return router, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}

	db, err := getDatabaseConnection()
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&ShortURL{}); err != nil {
		log.Fatalln("Failed to auto-migrate database:", err)
	}

	router, err := getGinRouter()
	if err != nil {
		log.Fatalln("Failed to create router:", err)
	}

	router.POST("/shorten", shortenURLHandler(db))
	router.GET("/:shortURL", redirectHandler(db))

	if err := router.Run(":8080"); err != nil {
		log.Fatalln("Failed to run server:", err)
	}
}
