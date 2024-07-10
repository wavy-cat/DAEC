package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wavy-cat/DAEC/backend/internal/config"
	"github.com/wavy-cat/DAEC/backend/internal/database"
	"github.com/wavy-cat/DAEC/backend/internal/utils/responses"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// AuthMiddleware - Middleware для авторизации запроса пользователя по JWT
type AuthMiddleware struct {
	Database *sql.DB      // Указатель на объект базы данных
	Logger   *zap.Logger  // Указатель на объект логгера
	Next     http.Handler // Функция, вызываемая middleware, которая будет обрабатывать http запрос
}

func (mw *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Получаем токен из запроса
	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		err := responses.RespondWithDefaultError(w, http.StatusUnauthorized)
		if err != nil {
			mw.Logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}

	// Проверяем его
	tokenString := strings.TrimPrefix(bearer, "Bearer ")
	tokenFromString, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := responses.RespondWithErrorMessage(w, http.StatusUnauthorized,
				fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
			if err != nil {
				mw.Logger.Error("failed to send response", zap.String("error", err.Error()))
			}
		}

		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		mw.Logger.Error("failed to verify JWT token", zap.String("error", err.Error()))
		err := responses.RespondWithPayload(w, http.StatusUnauthorized, err.Error())
		if err != nil {
			mw.Logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}

	claims, ok := tokenFromString.Claims.(jwt.MapClaims)
	if !ok {
		mw.Logger.Error("error when casting jwt.Claims to jwt.MapClaims")
		err := responses.RespondWithDefaultError(w, http.StatusInternalServerError)
		if err != nil {
			mw.Logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}

	// Получаем пользователя из базы данных
	user, err := database.SelectUserByLogin(context.TODO(), mw.Database, claims["name"].(string))
	if err != nil {
		mw.Logger.Error("failed to get user from database", zap.String("error", err.Error()))
		err := responses.RespondWithPayload(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			mw.Logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), "user", user))
	mw.Next.ServeHTTP(w, r)
}
