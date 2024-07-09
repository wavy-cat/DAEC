package utils

import "github.com/wavy-cat/DAEC/backend/internal/database"

type Expression struct {
	Id      int64   `json:"id"`
	Status  string  `json:"status"`  // Допустимые значения: pending, error, done
	Result  float64 `json:"result"`  // Пустое поле, в случае ошибки.
	Content string  `json:"content"` // Собственно само выражение
}

func ParseFromDBTypes(expressions []database.Expression) []Expression {
	exps := make([]Expression, 0, len(expressions))

	for _, exp := range expressions {
		exps = append(exps, Expression{
			Id:      exp.Id,
			Status:  exp.Status,
			Result:  exp.Result,
			Content: exp.Content,
		})
	}
	return exps
}

func ParseFromDBType(expression database.Expression) Expression {
	return Expression{
		Id:      expression.Id,
		Status:  expression.Status,
		Result:  expression.Result,
		Content: expression.Content,
	}
}
