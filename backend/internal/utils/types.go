package utils

import "github.com/google/uuid"

type ExpressionData struct {
	Id     uuid.UUID `json:"id"`
	Status string    `json:"status"` // Допустимые значения: pending, error, done
	Result float64   `json:"result"` // Пустое поле, в случае ошибки.
}

type ExpressionsSlice struct {
	Expressions []ExpressionData `json:"expressions"`
}

// Expression - Wrapper ExpressionData
type Expression struct {
	Expression ExpressionData `json:"expression"`
}
