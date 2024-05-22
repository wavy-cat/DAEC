package handler

import (
	"backend/internal/helper"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

type Data struct {
	Id         int    `json:"id"`         // Уникальный идентификатор выражения.
	Expression string `json:"expression"` // Строка с выражением.
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	// Получение логгера
	logger, ok := r.Context().Value("logger").(*zap.Logger)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Парсинг тела запроса
	data, err := helper.ParseJSON[Data](r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// TODO: Добавление в очередь
	logger.Info(fmt.Sprintf("%d %s", data.Id, data.Expression))

	// Отправка ответа
	_, err = fmt.Fprintf(w, "{\n}")
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
