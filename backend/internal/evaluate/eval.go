package evaluate

import (
	"context"
	"database/sql"
	"github.com/wavy-cat/DAEC/backend/internal/database"
	"github.com/wavy-cat/DAEC/backend/internal/tasks"
	"github.com/wavy-cat/DAEC/backend/pkg/postfix"
	"time"
)

// Evaluate вычисляет значение выражения в постфиксной записи.
// Ничего не возвращает, работает с БД напрямую.
func Evaluate(postfixNotation []any, id int64, db *sql.DB, manager *tasks.Manager) error {
	result := postfix.Calculate(postfixNotation, &solver{manager})

	for !result.IsDone {
		time.Sleep(100 * time.Millisecond)
	}

	var status string
	switch result.Error {
	case nil:
		status = "done"
	default:
		status = "error"
	}

	err := database.UpdateExpression(context.TODO(), db, &database.Expression{
		Id:     id,
		Status: status,
		Result: result.Result,
	})
	if err != nil {
		return err
	}

	return nil
}
