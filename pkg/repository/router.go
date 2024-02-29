package repository

import (
	"github.com/gin-gonic/gin"
	"os"
)

// GetGinRouter Функция для получения маршрутизатора Gin с установленными прокси и режимом работы из переменных окружения.
func GetGinRouter() (*gin.Engine, error) {
	gin.SetMode(os.Getenv("GIN_MODE"))
	router := gin.Default()
	err := router.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		return nil, err
	}
	return router, nil
}
