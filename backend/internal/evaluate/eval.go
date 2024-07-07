package evaluate

import (
	"github.com/google/uuid"
	"github.com/wavy-cat/DAEC/backend/internal/storage"
	"github.com/wavy-cat/DAEC/backend/internal/tasks"
	"github.com/wavy-cat/DAEC/backend/internal/utils"
	"github.com/wavy-cat/DAEC/backend/pkg/postfix"
	"time"
)

// Evaluate Вычисляет значение выражения в постфиксной записи.
// Ничего не возвращает, работает с БД напрямую.
// Необходимо запускать в другой горутине.
func Evaluate(postfixNotation []any, id uuid.UUID, storage *storage.Storage[utils.Expression], manager *tasks.Manager) {
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

	storage.Set(id, utils.Expression{
		Id:     id,
		Status: status,
		Result: result.Result,
	})
}
