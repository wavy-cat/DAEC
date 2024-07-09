package evaluate

import (
	"context"
	"database/sql"
	"github.com/wavy-cat/DAEC/backend/internal/database"
	"github.com/wavy-cat/DAEC/backend/internal/tasks"
	"github.com/wavy-cat/DAEC/backend/pkg/postfix"
)

// Recovery отправляет выражения в состоянии `pending` заново на обработку.
func Recovery(db *sql.DB, manager *tasks.Manager) error {
	// Берём все выражения из БД
	rawExps, err := database.SelectExpressions(context.TODO(), db)
	if err != nil {
		return err
	}

	// Фильтруем по статусу задачи
	exps := make([]database.Expression, 0)

	for _, val := range rawExps {
		if val.Status == "pending" {
			exps = append(exps, val)
		}
	}

	// Запускам обработку
	for _, task := range exps {
		go func(task *database.Expression) {
			notation, err := postfix.Convertor(task.Content)
			if err != nil {
				return
			}

			Evaluate(notation, task.Id, db, manager)
		}(&task)
	}

	return nil
}
