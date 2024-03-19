package main

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
)

const digits = "0123456789"
const lowercase = "abcdefghijklmnopqrstuvwxyz"
const uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const base62Digits = uppercase + lowercase + digits

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

func getStatusError(err error) int {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}
