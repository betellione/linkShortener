package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type ShortURL struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"` // Автоинкрементный ID
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
	tx.Model(base).Update("ShortURL", base.ShortURL)
	return nil
}

func findOrCreateURL(db *gorm.DB, inputURL string) (string, error) {
	var urlRecord ShortURL
	if err := db.Where("url = ?", inputURL).First(&urlRecord).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		urlRecord.URL = inputURL
		if err := db.Create(&urlRecord).Error; err != nil {
			return "", err
		}
	}
	return urlRecord.ShortURL, nil
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	dsn := "host=localhost user=postgres dbname=shorten password=q6J-LIFa6t port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	err = db.AutoMigrate(&ShortURL{})
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	router := gin.Default()
	router.POST("/shorten", func(context *gin.Context) {
		var jsonReq JsonRequest
		err = context.BindJSON(&jsonReq)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		shortened, err := findOrCreateURL(db, jsonReq.URL)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusOK, gin.H{"shortURL": shortened})
	})

	router.GET(":shortURL", func(context *gin.Context) {
		shortURL := context.Param("shortURL")
		var urlRecord ShortURL
		if err = db.Where("short_url = ?", shortURL).First(&urlRecord).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				context.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			} else {
				context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		context.Redirect(http.StatusFound, urlRecord.URL)
	})

	err = router.Run(":8084")
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
}
