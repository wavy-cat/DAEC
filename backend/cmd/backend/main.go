package main

import (
	"backend/http/handler"
	"backend/http/middleware"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

// MiddlewareOptions существует, для того чтобы складывать туда все аргументы для мидлварей.
// Данные поля будут доступны в контексте обработчика.
type MiddlewareOptions struct {
	Logger   *zap.Logger
	Database *map[string]interface{}
	Mutex    *sync.RWMutex
}

// setupMiddlewares оборачивает http.HandlerFunc в мидлвари.
// Используется в setupRouter
func setupMiddlewares(handlerFunc http.HandlerFunc, options MiddlewareOptions) http.Handler {
	// LoggingMiddleware -> DatabaseMiddleware -> HandlerFunc
	return &middleware.LoggingMiddleware{
		Logger: options.Logger,
		Next: &middleware.DatabaseMiddleware{
			Database: options.Database,
			Lock:     options.Mutex,
			Next:     handlerFunc,
		},
	}
}

// setupRouter создаёт новый экземпляр mux.Router и настраивает конфигурацию маршрутизации.
func setupRouter(options MiddlewareOptions) *mux.Router {
	router := mux.NewRouter()

	// Здесь настраивается маршрутизация.
	// Чтобы не оборачивать функции в мидлвари вручную используется setupMiddlewares.
	router.Handle("/api/v1/calculate", setupMiddlewares(handler.CalculateHandler, options)).Methods("POST")

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

	// Создание базы данных и других опций для MiddlewareOptions
	database := make(map[string]interface{})
	options := MiddlewareOptions{
		Logger:   logger,
		Database: &database,
		Mutex:    &sync.RWMutex{},
	}
	// TODO: Добавить Tasks{}

	if err := http.ListenAndServe(":8080", setupRouter(options)); err != nil {
		logger.Fatal(err.Error())
		return
	}
}
