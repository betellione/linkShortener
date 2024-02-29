package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ShortURL Определение структуры для записи коротких URL в базу данных.
type ShortURL struct {
	gorm.Model
	ShortURL string `gorm:"unique_index"`
	URL      string `gorm:"unique_index"`
}

// JsonRequest Структура для парсинга JSON запроса с полным URL.
type JsonRequest struct {
	URL string `json:"url"`
}

const base62Digits = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// To62Base Функция конвертации числового ID в короткий код на основе 62 символов.
func To62Base(num uint) string {
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

// AfterCreate Хук, автоматически срабатывающий после создания записи в базе данных для генерации короткого URL.
func (base *ShortURL) AfterCreate(tx *gorm.DB) (err error) {
	base.ShortURL = To62Base(base.ID)
	return tx.Model(base).Update("ShortURL", base.ShortURL).Error
}

// Функция поиска существующего или создания нового короткого URL.
func findOrCreateURL(db *gorm.DB, inputURL string) (string, error) {
	urlRecord := ShortURL{URL: inputURL}
	return urlRecord.ShortURL, db.Where("url = ?", inputURL).FirstOrCreate(&urlRecord).Error
}

// Функция для подключения к базе данных с параметрами из переменных окружения.
func getDatabaseConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("DB_PORT"))
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// Определение обработчика HTTP POST запросов для создания коротких URL.
func postHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var jsonReq JsonRequest
		if err := c.BindJSON(&jsonReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		shortened, err := findOrCreateURL(db, jsonReq.URL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or find URL"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"shortId": shortened})
	}
}

// Определение обработчика HTTP GET запросов для перенаправления по короткому URL.
func getHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortURL := c.Param("shortURL")
		var urlRecord ShortURL
		if err := db.Where("short_url = ?", shortURL).First(&urlRecord).Error; err != nil {
			c.JSON(getStatusError(err), gin.H{"error": "URL not found"})
			return
		}
		c.Redirect(http.StatusFound, urlRecord.URL)
	}
}

// Вспомогательная функция для определения статуса ошибки при общении с базой данных.
func getStatusError(err error) (status int) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

// Функция для получения маршрутизатора Gin с установленными прокси и режимом работы из переменных окружения.
func getGinRouter() (*gin.Engine, error) {
	gin.SetMode(os.Getenv("GIN_MODE"))
	router := gin.Default()
	err := router.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		return nil, err
	}
	return router, nil
}

// Главная функция программы для загрузки .env файла, подключения к базе данных,
// миграции схемы и настройки HTTP сервера.
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
		log.Fatalln("Failed to initialize Gin router:", err)
	}

	router.POST("/create", postHandler(db))
	router.GET("/:shortURL", getHandler(db))

	if err := router.Run(); err != nil {
		log.Fatalln("Failed to run the server:", err)
	}
}
