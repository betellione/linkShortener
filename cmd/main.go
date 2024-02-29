package main

import (
	"log"
	"shorten/pkg/handler"
	"shorten/pkg/repository"
	"shorten/pkg/structs"

	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// @title Link Shortener
// @version 1.0
// @description API Server for TodoList Application

// @host localhost:8080
// @BasePath /
func main() {

	// TODO сваггер
	// TODO cashing
	// TODO web view
	// TODO CI/CD
	// TODO Viper
	// TODO Allure

	db, err := repository.GetDatabaseConnection()
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&structs.ShortURL{}); err != nil {
		log.Fatalln("Failed to auto-migrate database:", err)
	}

	router, err := repository.GetGinRouter()
	if err != nil {
		log.Fatalln("Failed to initialize Gin router:", err)
	}

	router.POST("/create", handler.PostHandler(db))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//router.GET("/:shortURL", handler.RedirectHandler(db))

	if err := router.Run(); err != nil {
		log.Fatalln("Failed to run the server:", err)
	}
}
