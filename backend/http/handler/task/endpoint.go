package task

import (
	"backend/internal/tasks"
	"backend/internal/utils"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		HandlerPost(w, r)
		return
	}

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

	// Получение задачи
	task, found := manager.GetTask()
	if !found {
		err := utils.RespondWithDefaultError(w, http.StatusNotFound)
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	err := utils.RespondWithPayload(w, http.StatusOK, task)
	if err != nil {
		logger.Error(err.Error())
	}
}

func HandlerPost(w http.ResponseWriter, r *http.Request) {
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
	var data tasks.ResultData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		err := utils.RespondWithErrorMessage(w, http.StatusUnprocessableEntity, err.Error())
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	// Отправка результата
	fmt.Println(data.Result)
	err = manager.AddResultToTask(data.Id, data.Result.Float, data.Result.Valid)
	if err != nil {
		err := utils.RespondWithErrorMessage(w, http.StatusNotFound, err.Error())
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	utils.RespondOnlyCode(w, http.StatusOK)
}
