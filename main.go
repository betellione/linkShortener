package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

// JsonRequest Определение структур
type JsonRequest struct {
	URL string `json:"url"`
}

type ShortURL struct {
	gorm.Model
	ShortURL string `gorm:"unique_index"`
	URL      string `gorm:"unique_index"`
}

// Database Интерфейс Database
type Database interface {
	FindOrCreateURL(inputURL string) (string, error)
	GetURL(shortURL string) (string, error)
}

type GormDatabase struct {
	db *gorm.DB
}

const digits = "0123456789"
const lowercase = "abcdefghijklmnopqrstuvwxyz"
const uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const base62Digits = uppercase + lowercase + digits

// NewGormDatabase Инициализация новой БД
func NewGormDatabase(dsn string) (*GormDatabase, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &GormDatabase{db: db}, nil
}

// FindOrCreateURL Реализация интерфейса Database для GormDatabase
func (gdb *GormDatabase) FindOrCreateURL(inputURL string) (string, error) {
	urlRecord := ShortURL{URL: inputURL}
	if err := gdb.db.Where("url = ?", inputURL).FirstOrCreate(&urlRecord).Error; err != nil {
		return "", err
	}
	return urlRecord.ShortURL, nil
}

func (gdb *GormDatabase) GetURL(shortURL string) (string, error) {
	var urlRecord ShortURL
	if err := gdb.db.Where("short_url = ?", shortURL).First(&urlRecord).Error; err != nil {
		return "", err
	}
	return urlRecord.URL, nil
}

// AfterCreate Методы ShortURL
func (base *ShortURL) AfterCreate(tx *gorm.DB) (err error) {
	base.ShortURL = to62Base(base.ID)
	return tx.Model(base).Update("ShortURL", base.ShortURL).Error
}

// PostHandler Функции маршрутизации
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

// Вспомогательные функции
func getStatusError(err error) int {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func to62Base(num uint) string {
	if num == 0 {
		return "0"
	}
	base := uint(62)
	result := ""
	for num > 0 {
		remainder := num % base
		num /= base
		result = string(base62Digits[remainder]) + result
	}
	return result
}

func main() {
	// Инициализация базы данных и маршрутизатора Gin
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("DB_PORT"))
	db, err := NewGormDatabase(dsn)
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}

	// Автомиграция для ShortURL
	if err := db.db.AutoMigrate(&ShortURL{}); err != nil {
		log.Fatalln("Failed to auto-migrate database:", err)
	}

	// Инициализация Gin роутера
	router, err := getGinRouter()
	if err != nil {
		log.Fatalln("Failed to initialize Gin router:", err)
	}
	router.POST("/create", PostHandler(db))
	router.GET("/:shortURL", RedirectHandler(db))

	if err := router.Run(); err != nil {
		log.Fatalln("Failed to run the server:", err)
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
