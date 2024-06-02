package evaluate

import (
	"backend/internal/storage"
	"backend/internal/tasks"
	"backend/internal/utils"
	"backend/pkg/postfix"
	"github.com/google/uuid"
	"time"
)

// Evaluate Вычисляет значение выражения в постфиксной записи.
// Ничего не возвращает, работает с БД напрямую.
// Необходимо запускать в другой горутине.
func Evaluate(postfixNotation []any, id uuid.UUID, storage *storage.Storage[utils.ExpressionData], manager *tasks.Manager) {
	result := postfix.Calculate(postfixNotation, &solver{manager})

	for !result.IsDone {
		time.Sleep(100 * time.Millisecond)
	}

	var status string
	if result.Error != nil {
		status = "error"
	} else {
		status = "done"
	}

	storage.Set(id, utils.ExpressionData{
		Id:     id,
		Status: status,
		Result: result.Result,
	})
}
