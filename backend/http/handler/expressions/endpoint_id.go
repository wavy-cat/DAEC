package expressions

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wavy-cat/DAEC/backend/internal/database"
	"github.com/wavy-cat/DAEC/backend/internal/utils"
	"github.com/wavy-cat/DAEC/backend/internal/utils/responses"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func HandlerById(w http.ResponseWriter, r *http.Request) {
	// Получение логгера
	logger, ok := r.Context().Value("logger").(*zap.Logger)
	if !ok {
		fmt.Println("failed to get logger in expression by id")
		err := responses.RespondWithDefaultError(w, http.StatusInternalServerError)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	// Получение пользователя
	user, ok := r.Context().Value("user").(database.User)
	if !ok {
		logger.Error("failed to get user")
		err := responses.RespondWithDefaultError(w, http.StatusInternalServerError)
		if err != nil {
			logger.Error("failed to send response", zap.String("error", err.Error()))
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

	rawId, ok := mux.Vars(r)["id"]
	if !ok {
		// К сожалению, по условию задачи должны быть коды 200, 404 или 500.
		// Но в данном случаю лучше было бы выдать 400 Bad Request.
		err := responses.RespondWithErrorMessage(w, http.StatusNotFound, "id not sent")
		if err != nil {
			logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}

	id, err := strconv.Atoi(rawId)
	if err != nil {
		err := responses.RespondWithErrorMessage(w, http.StatusUnprocessableEntity, err.Error())
		if err != nil {
			logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}

	exp, err := database.SelectUserExpressionByID(context.TODO(), db, user.Id, int64(id))
	if err != nil {
		err := responses.RespondWithDefaultError(w, http.StatusNotFound)
		if err != nil {
			logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}

	err = responses.RespondWithPayload(
		w, http.StatusOK,
		struct {
			Expression utils.Expression `json:"expression"`
		}{Expression: utils.ParseFromDBType(exp)})
	if err != nil {
		logger.Error("failed to send response", zap.String("error", err.Error()))
	}
}
