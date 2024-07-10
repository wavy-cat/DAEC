package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/wavy-cat/DAEC/backend/internal/database"
	"github.com/wavy-cat/DAEC/backend/internal/utils/responses"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Получение логгера
	logger, ok := r.Context().Value("logger").(*zap.Logger)
	if !ok {
		fmt.Println("failed to get logger in register")
		err := responses.RespondWithDefaultError(w, http.StatusInternalServerError)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	// Получение бд
	db, ok := r.Context().Value("database").(*sql.DB)
	if !ok {
		logger.Error("failed to get database")
		err := responses.RespondWithDefaultError(w, http.StatusInternalServerError)
		if err != nil {
			logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}

	// Парсинг тела запроса
	var data Request
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		if err := responses.RespondWithErrorMessage(w, http.StatusUnprocessableEntity, err.Error()); err != nil {
			logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}

	// Хэшируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed to create password hash", zap.String("error", err.Error()))

		if err := responses.RespondWithDefaultError(w, http.StatusInternalServerError); err != nil {
			logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}

	// Добавляем нового пользователя в базу данных
	id, err := database.InsertUser(context.TODO(), db, &database.User{Login: data.Login, Password: hash})
	if err != nil {
		switch err.Error() {
		case "UNIQUE constraint failed: users.login":
			err = responses.RespondWithErrorMessage(w, http.StatusConflict, "a user with this login is already registered")
		default:
			logger.Error("failed to add new user", zap.String("error", err.Error()))
			err = responses.RespondWithDefaultError(w, http.StatusInternalServerError)
		}

		if err != nil {
			logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}

	logger.Info("new user added", zap.Int64("id", id))

	responses.RespondOnlyCode(w, http.StatusOK)
}
