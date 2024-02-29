package repository

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

// GetDatabaseConnection Функция для подключения к базе данных с параметрами из переменных окружения.
func GetDatabaseConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("DB_PORT"))
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
