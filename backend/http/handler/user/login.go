package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wavy-cat/DAEC/backend/internal/config"
	"github.com/wavy-cat/DAEC/backend/internal/database"
	"github.com/wavy-cat/DAEC/backend/internal/utils/responses"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Получение логгера
	logger, ok := r.Context().Value("logger").(*zap.Logger)
	if !ok {
		fmt.Println("failed to get logger in login")
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
		err := responses.RespondWithErrorMessage(w, http.StatusUnprocessableEntity, err.Error())
		if err != nil {
			logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}

	// Получаем пользователя из базы данных
	user, err := database.SelectUserByLogin(context.TODO(), db, data.Login)
	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			err = responses.RespondWithErrorMessage(w, http.StatusNotFound, "user is not found")
		default:
			logger.Error("failed to get user from database", zap.String("error", err.Error()))
			err = responses.RespondWithDefaultError(w, http.StatusInternalServerError)
		}

		if err != nil {
			logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(data.Password))
	if err != nil {
		// Пароли не совпадают
		err := responses.RespondWithErrorMessage(w, http.StatusForbidden, "incorrect password entered")
		if err != nil {
			logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}

	// Генерируем JWT токен
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": user.Login,
		"nbf":  now.Unix(),
		"exp":  now.Add(time.Hour).Unix(),
		"iat":  now.Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		err := responses.RespondWithDefaultError(w, http.StatusInternalServerError)
		if err != nil {
			logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}

	err = responses.RespondWithPayload(w, http.StatusOK, LoginResponse{
		Token:  tokenString,
		Expire: now.Add(time.Hour).Unix(),
	})
	if err != nil {
		logger.Error("failed to send response", zap.String("error", err.Error()))
	}
}
