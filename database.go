package main

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	FindOrCreateURL(inputURL string) (string, error)
	GetURL(shortURL string) (string, error)
}

type GormDatabase struct {
	db *gorm.DB
}

func NewGormDatabase() (*GormDatabase, error) {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable",
		viper.GetString("db_host"),
		viper.GetString("postgres_user"),
		viper.GetString("postgres_db"),
		viper.GetString("postgres_password"),
		viper.GetString("db_port"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &GormDatabase{db: db}, nil
}

// AfterCreate Методы ShortURL
func (base *ShortURL) AfterCreate(tx *gorm.DB) (err error) {
	base.ShortURL = to62Base(base.ID)
	return tx.Model(base).Update("ShortURL", base.ShortURL).Error
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
