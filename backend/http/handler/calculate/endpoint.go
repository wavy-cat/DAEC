package calculate

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/wavy-cat/DAEC/backend/internal/database"
	"github.com/wavy-cat/DAEC/backend/internal/evaluate"
	"github.com/wavy-cat/DAEC/backend/internal/tasks"
	"github.com/wavy-cat/DAEC/backend/internal/utils"
	"github.com/wavy-cat/DAEC/backend/pkg/postfix"
	"go.uber.org/zap"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Получение логгера
	logger, ok := r.Context().Value("logger").(*zap.Logger)
	if !ok {
		fmt.Println("failed to get logger in calculate")
		err := utils.RespondWithDefaultError(w, http.StatusInternalServerError)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	// Получение бд
	db, ok := r.Context().Value("database").(*sql.DB)
	if !ok {
		logger.Error("failed to get database")
		err := utils.RespondWithDefaultError(w, http.StatusInternalServerError)
		if err != nil {
			logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}

	// Получение менеджера задач
	manager, ok := r.Context().Value("manager").(*tasks.Manager)
	if !ok {
		logger.Error("failed to get tasks manager")
		err := utils.RespondWithDefaultError(w, http.StatusInternalServerError)
		if err != nil {
			logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}

	// Парсинг тела запроса
	var data DataRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		err := utils.RespondWithErrorMessage(w, http.StatusUnprocessableEntity, err.Error())
		if err != nil {
			logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}
	expression := data.Expression

	// Проверка, что выражение не пустое
	if expression == "" {
		err := utils.RespondWithErrorMessage(w, http.StatusUnprocessableEntity, "expression must not be empty")
		if err != nil {
			logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}

	// Проверка, что выражение не содержит запрещённые символы
	var allowedChars = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '+', '-', '(', ')', '/', '*', '^', '.', ' '}
	if contains, char := utils.CheckCharsInString(expression, allowedChars); !contains {
		err := utils.RespondWithErrorMessage(w, http.StatusUnprocessableEntity, "forbidden symbol found: "+string(char))
		if err != nil {
			logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}

	// Переводим выражение в постфиксную запись (пытаемся)
	postfixNotation, err := postfix.Convertor(expression)
	if err != nil {
		// Какая-то проблема с выражением
		err := utils.RespondWithErrorMessage(w, http.StatusUnprocessableEntity, err.Error())
		if err != nil {
			logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}

	// Добавляем выражение в БД
	id, err := database.InsertExpression(context.TODO(), db, &database.Expression{Status: "pending", Content: expression})
	if err != nil {
		err := utils.RespondWithDefaultError(w, http.StatusInternalServerError)
		if err != nil {
			logger.Error("failed to add expression to database", zap.String("error", err.Error()))
		}
		return
	}

	// Отправка ответа
	err = utils.RespondWithPayload(w, http.StatusCreated, DataResponse{Id: id})
	if err != nil {
		logger.Error("failed to send response", zap.String("error", err.Error()))
	}

	// Отправляем на обработку
	go func(postfixNotation []any, id int64, db *sql.DB, manager *tasks.Manager) {
		err := evaluate.Evaluate(postfixNotation, id, db, manager)
		if err != nil {
			logger.Error("error from Evaluate", zap.String("error", err.Error()))
		}
	}(postfixNotation, id, db, manager)
}
