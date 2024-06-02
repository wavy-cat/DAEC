package expressions

import (
	stg "backend/internal/storage"
	"backend/internal/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
	storage, ok := r.Context().Value("storage").(*stg.Storage[utils.ExpressionData])
	if !ok {
		logger.Error("Failed to get storage")
		err := utils.RespondWithDefaultError(w, http.StatusInternalServerError)
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	exps := storage.GetAll()
	err := utils.RespondWithPayload(w, http.StatusOK, utils.ExpressionsSlice{Expressions: exps})
	if err != nil {
		logger.Error(err.Error())
	}
}

func HandlerById(w http.ResponseWriter, r *http.Request) {
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
	storage, ok := r.Context().Value("storage").(*stg.Storage[utils.ExpressionData])
	if !ok {
		logger.Error("Failed to get storage")
		err := utils.RespondWithDefaultError(w, http.StatusInternalServerError)
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	rawId := mux.Vars(r)["id"]
	id, err := uuid.Parse(rawId)
	if err != nil {
		err := utils.RespondWithErrorMessage(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	exp, found := storage.Get(id)
	if !found {
		err := utils.RespondWithDefaultError(w, http.StatusNotFound)
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	err = utils.RespondWithPayload(w, http.StatusOK, utils.Expression{Expression: exp})
	if err != nil {
		logger.Error(err.Error())
	}
}
