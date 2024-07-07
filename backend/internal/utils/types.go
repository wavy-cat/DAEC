package utils

import "github.com/google/uuid"

type Expression struct {
	Id     uuid.UUID `json:"id"`
	Status string    `json:"status"` // Допустимые значения: pending, error, done
	Result float64   `json:"result"` // Пустое поле, в случае ошибки.
}
