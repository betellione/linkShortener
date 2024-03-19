package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExample(t *testing.T) {
	actualResult := "ожидаемое значение"
	expectedResult := "ожидаемое значение"

	assert.Equal(t, expectedResult, actualResult, "Результат функции должен соответствовать ожидаемому")
}
