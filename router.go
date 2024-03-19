package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

func getGinRouter() (*gin.Engine, error) {
	gin.SetMode(os.Getenv("GIN_MODE"))
	router := gin.Default()
	err := router.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		return nil, err
	}
	return router, nil
}
