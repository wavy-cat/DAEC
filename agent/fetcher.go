package main

import (
	"agent/config"
	"agent/work"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// Task описывает задачу, которую необходимо выполнить
type Task struct {
	Id            uuid.UUID `json:"id"`             // Уникальный идентификатор задачи (UUID)
	Arg1          float64   `json:"arg1"`           // Первый аргумент операции
	Arg2          float64   `json:"arg2"`           // Второй аргумент операции
	Operation     string    `json:"operation"`      // Операция, которую следует выполнить
	OperationTime int       `json:"operation_time"` // Время, требуемое для выполнения операции в миллисекундах
}

// TaskWrapper представляет собой обёртку вокруг задач, она используется для корректного преобразования вложенных JSON.
// Его использование обязательно.
type TaskWrapper struct {
	Task Task `json:"task"` // Task содержит в себе структуру задачи
}

// Функция для получения новых задач от оркестратора
func fetcher(pool *work.Pool, logger *zap.Logger) {
	for {
		// Отправляем запрос оркестратору
		response, err := http.Get(config.BackendUrl + "/internal/task")
		if err != nil {
			logger.Error("Failed to send request to orchestrator. Try again in 2 seconds...")
			time.Sleep(2 * time.Second)
			continue
		}

		// Проверяем ответ
		switch response.StatusCode {
		case 404:
			// Если новой задачи нет, то ненадолго "засыпаем"
			time.Sleep(300 * time.Millisecond)
			continue
		case 500:
			logger.Error("An internal orchestrator error occurred. Try again in 2 seconds...")
			time.Sleep(2 * time.Second)
			continue
		case 200:
			// Скипаем эту часть, она реализована ниже
		default:
			logger.Error(fmt.Sprintf("Unexpected response from the server (%d). Try again in 2 seconds...", response.StatusCode))
			time.Sleep(2 * time.Second)
			continue
		}

		// Демаршализируем ответ
		var taskObj TaskWrapper
		err = json.NewDecoder(response.Body).Decode(&taskObj)
		if err != nil {
			logger.Error("Error unmarshalling JSON in request. Task skipped")
			continue
		}

		// Создаём новую задачу и отправляем её в пул
		task := taskObj.Task
		expression := work.Expression{
			Id:            task.Id,
			Num1:          task.Arg1,
			Num2:          task.Arg2,
			Operator:      rune(task.Operation[0]),
			OperationTime: time.Duration(task.OperationTime) * time.Millisecond,
		}
		pool.Run(expression)
	}
}
