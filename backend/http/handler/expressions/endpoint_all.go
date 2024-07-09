package expressions

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/wavy-cat/DAEC/backend/internal/database"
	"github.com/wavy-cat/DAEC/backend/internal/utils"
	"go.uber.org/zap"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Получение логгера
	logger, ok := r.Context().Value("logger").(*zap.Logger)
	if !ok {
		fmt.Println("failed to get logger in expressions")
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

	exps, err := database.SelectExpressions(context.TODO(), db)
	if err != nil {
		logger.Error("failed to expressions from database", zap.String("error", err.Error()))
		err := utils.RespondWithDefaultError(w, http.StatusInternalServerError)
		if err != nil {
			logger.Error("failed to send response", zap.String("error", err.Error()))
		}
		return
	}

	err = utils.RespondWithPayload(
		w, http.StatusOK,
		struct {
			Expressions []utils.Expression `json:"expressions"`
		}{Expressions: utils.ParseFromDBTypes(exps)})
	if err != nil {
		logger.Error("failed to send response", zap.String("error", err.Error()))
	}
}
