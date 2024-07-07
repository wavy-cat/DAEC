package main

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/wavy-cat/DAEC/backend/http/handler/calculate"
	"github.com/wavy-cat/DAEC/backend/http/handler/expressions"
	"github.com/wavy-cat/DAEC/backend/http/middleware"
	"github.com/wavy-cat/DAEC/backend/internal/config"
	"github.com/wavy-cat/DAEC/backend/internal/storage"
	"github.com/wavy-cat/DAEC/backend/internal/tasks"
	"github.com/wavy-cat/DAEC/backend/internal/utils"
	"go.uber.org/zap"
	"net/http"
)

// setupMiddlewares оборачивает http.HandlerFunc в мидлвари.
// Используется в setupRouter
func setupMiddlewares(handlerFunc http.HandlerFunc,
	logger *zap.Logger,
	storage *storage.Storage[utils.Expression],
	manager *tasks.Manager) http.Handler {
	// DatabaseMiddleware -> ManagerMiddleware -> LoggingMiddleware -> HandlerFunc
	return &middleware.DatabaseMiddleware{
		Storage: storage,
		Next: &middleware.ManagerMiddleware{
			Manager: manager,
			Next: &middleware.LoggingMiddleware{
				Logger: logger,
				Next:   handlerFunc,
			},
		},
	}
}

// setupRouter создаёт новый экземпляр mux.Router и настраивает конфигурацию маршрутизации.
func setupRouter(logger *zap.Logger,
	storage *storage.Storage[utils.Expression],
	manager *tasks.Manager) http.Handler {
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
	}

	for path, routeConfig := range routes {
		// Чтобы не оборачивать функции в мидлвари вручную используется setupMiddlewares
		handler := setupMiddlewares(routeConfig.handler, logger, storage, manager)
		router.Handle(path, handler).Methods(routeConfig.methods...)
	}

	return cors.Default().Handler(router)
}

func startHTTPServer(logger *zap.Logger, db *storage.Storage[utils.Expression], manager *tasks.Manager) {
	logger.Info("Starting the HTTP server...")
	if err := http.ListenAndServe(config.HTTPAddress, setupRouter(logger, db, manager)); err != nil {
		logger.Fatal(err.Error())
	}
}
