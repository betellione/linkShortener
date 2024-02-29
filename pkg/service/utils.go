package service

import (
	"gorm.io/gorm"
	"shorten/pkg/structs"
)

// FindOrCreateURL Функция поиска существующего или создания нового короткого URL.
func FindOrCreateURL(db *gorm.DB, inputURL string) (string, error) {
	urlRecord := structs.ShortURL{URL: inputURL}
	return urlRecord.ShortURL, db.Where("url = ?", inputURL).FirstOrCreate(&urlRecord).Error
}
