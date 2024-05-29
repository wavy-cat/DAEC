package main

import (
	"backend/http/handler/calculate"
	"backend/http/middleware"
	"backend/internal/storage"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

// setupMiddlewares оборачивает http.HandlerFunc в мидлвари.
// Используется в setupRouter
func setupMiddlewares(handlerFunc http.HandlerFunc, logger *zap.Logger, storage *storage.Storage) http.Handler {
	// LoggingMiddleware -> DatabaseMiddleware -> HandlerFunc
	return &middleware.LoggingMiddleware{
		Logger: logger,
		Next: &middleware.DatabaseMiddleware{
			Storage: storage,
			Next:    handlerFunc,
		},
	}
}

// setupRouter создаёт новый экземпляр mux.Router и настраивает конфигурацию маршрутизации.
func setupRouter(logger *zap.Logger, storage *storage.Storage) *mux.Router {
	router := mux.NewRouter()

	// Здесь настраивается маршрутизация aka указание эндпойнтов сервера
	routes := map[string]struct {
		handler http.HandlerFunc
		methods []string
	}{
		// Указываются тут, если что.
		// Ключ — путь, значение — структура из объекта http.HandlerFunc и разрешённых методов.
		"/api/v1/calculate": {
			calculate.Handler,
			[]string{"POST"},
		},
	}

	for path, routeConfig := range routes {
		// Чтобы не оборачивать функции в мидлвари вручную используется setupMiddlewares
		router.Handle(path, setupMiddlewares(routeConfig.handler, logger, storage)).Methods(routeConfig.methods...)
	}

	return router
}

func main() {
	// Создание и запуск логгера
	logger, err := zap.NewDevelopment() // Заменить в конце на NewProduction

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

	// Создание базы данных
	db := storage.NewStorage()

	// Запуск сервера
	if err := http.ListenAndServe(":8080", setupRouter(logger, db)); err != nil {
		logger.Fatal(err.Error())
	}
}
