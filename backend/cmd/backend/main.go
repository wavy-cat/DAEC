package main

import (
	"backend/http/handler/calculate"
	"backend/http/handler/expressions"
	"backend/http/handler/task"
	"backend/http/middleware"
	"backend/internal/config"
	"backend/internal/storage"
	"backend/internal/tasks"
	"backend/internal/utils"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

// setupMiddlewares оборачивает http.HandlerFunc в мидлвари.
// Используется в setupRouter
func setupMiddlewares(handlerFunc http.HandlerFunc, logger *zap.Logger,
	storage *storage.Storage[utils.ExpressionData],
	manager *tasks.Manager) http.Handler {
	// LoggingMiddleware -> DatabaseMiddleware -> ManagerMiddleware -> HandlerFunc
	return &middleware.LoggingMiddleware{
		Logger: logger,
		Next: &middleware.DatabaseMiddleware{
			Storage: storage,
			Next: &middleware.ManagerMiddleware{
				Manager: manager,
				Next:    handlerFunc,
			},
		},
	}
}

// setupRouter создаёт новый экземпляр mux.Router и настраивает конфигурацию маршрутизации.
func setupRouter(logger *zap.Logger, storage *storage.Storage[utils.ExpressionData], manager *tasks.Manager) *mux.Router {
	router := mux.NewRouter()

	// Здесь настраивается маршрутизация aka указание эндпойнтов сервера
	routes := map[string]struct {
		handler http.HandlerFunc
		methods []string
	}{
		// Указываются тут, если что.
		// Ключ — путь, значение — структура из объекта HandlerFunc и разрешённых методов.
		"/api/v1/calculate": {
			calculate.Handler,
			[]string{"POST"},
		},
		"/api/v1/expressions": {
			expressions.Handler,
			[]string{"GET"},
		},
		"/api/v1/expressions/{id}": {
			expressions.HandlerById,
			[]string{"GET"},
		},
		"/internal/task": {
			task.Handler,
			[]string{"GET", "POST"},
		},
	}

	for path, routeConfig := range routes {
		// Чтобы не оборачивать функции в мидлвари вручную используется setupMiddlewares
		handler := setupMiddlewares(routeConfig.handler, logger, storage, manager)
		router.Handle(path, handler).Methods(routeConfig.methods...)
	}

	return router
}

func main() {
	// Создание и запуск логгера
	logger, err := zap.NewProduction()

	if err != nil {
		fmt.Println("error initializing logger:", err)
		return
	}

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Println("error synchronizing logger:", err)
		}
	}(logger)

	logger.Info("Orchestrator is starting")

	// Создание базы данных и очереди
	db := storage.NewStorage[utils.ExpressionData]()
	manager := tasks.NewManager()
	defer manager.ShutdownWatcher()

	// Запуск сервера
	if err := http.ListenAndServe(config.ServerAddress, setupRouter(logger, db, manager)); err != nil {
		logger.Fatal(err.Error())
	}
}
