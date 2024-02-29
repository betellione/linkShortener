package structs

import (
	"gorm.io/gorm"
)

// ShortURL Определение структуры для записи коротких URL в базу данных.
type ShortURL struct {
	gorm.Model
	ShortURL string `gorm:"unique_index"`
	URL      string `gorm:"unique_index"`
}

// AfterCreate Хук, автоматически срабатывающий после создания записи в базе данных для генерации короткого URL.
func (base *ShortURL) AfterCreate(tx *gorm.DB) (err error) {
	base.ShortURL = to62Base(base.ID)
	return tx.Model(base).Update("ShortURL", base.ShortURL).Error
}

const digits = "0123456789"
const lowercase = "abcdefghijklmnopqrstuvwxyz"
const uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const base62Digits = uppercase + lowercase + digits

// To62Base Функция конвертации числового ID в короткий код на основе 62 символов.
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
