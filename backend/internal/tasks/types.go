package tasks

import (
	"github.com/google/uuid"
	"time"
)

type TaskData struct {
	Id            uuid.UUID `json:"id"`             // Уникальный идентификатор задачи (UUID)
	Arg1          float64   `json:"arg1"`           // Первый аргумент операции
	Arg2          float64   `json:"arg2"`           // Второй аргумент операции
	Operation     byte      `json:"operation"`      // Операция, которую следует выполнить
	OperationTime uint32    `json:"operation_time"` // Время, требуемое для выполнения операции в миллисекундах
}

type Task struct {
	Data           TaskData
	Status         string // Возможные значения: queue, processing, done
	Result         float64
	Successful     bool // Успешно ли выполнена задача
	Timeout        time.Duration
	CompleteBefore time.Time // Время, до которого задача должна быть решена
}

type TaskResult struct {
	IsDone     bool
	Result     float64
	Successful bool
}
