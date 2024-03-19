package main

import "gorm.io/gorm"

type JsonRequest struct {
	URL string `json:"url"`
}

type ShortURL struct {
	gorm.Model
	ShortURL string `gorm:"unique_index"`
	URL      string `gorm:"unique_index"`
}
