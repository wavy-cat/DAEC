package calculate

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/wavy-cat/DAEC/backend/internal/evaluate"
	stg "github.com/wavy-cat/DAEC/backend/internal/storage"
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
		fmt.Println("Failed to get logger in calculate")
		err := utils.RespondWithDefaultError(w, http.StatusInternalServerError)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	// Получение бд
	storage, ok := r.Context().Value("storage").(*stg.Storage[utils.Expression])
	if !ok {
		logger.Error("Failed to get storage")
		err := utils.RespondWithDefaultError(w, http.StatusInternalServerError)
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	// Получение менеджера задач
	manager, ok := r.Context().Value("manager").(*tasks.Manager)
	if !ok {
		logger.Error("Failed to get tasks manager")
		err := utils.RespondWithDefaultError(w, http.StatusInternalServerError)
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	// Парсинг тела запроса
	var data DataRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		err := utils.RespondWithErrorMessage(w, http.StatusUnprocessableEntity, err.Error())
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	// Проверка, что выражение не пустое
	if data.Expression == "" {
		err := utils.RespondWithErrorMessage(w, http.StatusUnprocessableEntity, "expression must not be empty")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	// Проверка, что выражение не содержит запрещённые символы
	var allowedChars = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '+', '-', '(', ')', '/', '*', '^', '.', ' '}
	if contains, char := utils.CheckCharsInString(data.Expression, allowedChars); !contains {
		err := utils.RespondWithErrorMessage(w, http.StatusUnprocessableEntity, "forbidden symbol found: "+string(char))
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	// Переводим выражение в постфиксную запись (пытаемся)
	postfixNotation, err := postfix.Convertor(data.Expression)
	if err != nil {
		// Какая-то проблема с выражением
		err := utils.RespondWithErrorMessage(w, http.StatusUnprocessableEntity, err.Error())
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	// Генерируем ID
	id, err := uuid.NewRandom() // генерируем id
	if err != nil {
		logger.Error(err.Error())
		err := utils.RespondWithDefaultError(w, http.StatusInternalServerError)
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	// Отправка ответа
	err = utils.RespondWithPayload(w, http.StatusCreated, DataResponse{Id: id})
	if err != nil {
		logger.Error(err.Error())
	}
	err = r.Body.Close()
	if err != nil {
		logger.Error(err.Error())
	}

	// Добавляем запись в БД
	storage.Set(id, utils.Expression{Id: id, Status: "pending"})

	// Отправляем на обработку
	go evaluate.Evaluate(postfixNotation, id, storage, manager)
}
