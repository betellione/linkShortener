package main

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	_ "shorten/docs"
)

// @title Shortener API
// @description API для сокращения URL.
// @version 1.0
// @host localhost:8080
// @BasePath /
func main() {
	initConfig()
	db, err := NewGormDatabase()
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/create", PostHandler(db))
	router.GET("/:shortURL", RedirectHandler(db))

	if err := router.Run(); err != nil {
		log.Fatalln("Failed to run the server:", err)
	}
}
