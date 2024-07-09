package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/wavy-cat/DAEC/backend/http/handler/calculate"
	"github.com/wavy-cat/DAEC/backend/http/handler/expressions"
	"github.com/wavy-cat/DAEC/backend/http/middleware"
	"github.com/wavy-cat/DAEC/backend/internal/config"
	"github.com/wavy-cat/DAEC/backend/internal/tasks"
	"go.uber.org/zap"
	"net/http"
)

// setupMiddlewares оборачивает http.HandlerFunc в мидлвари.
// Используется в setupRouter
func setupMiddlewares(handlerFunc http.HandlerFunc, logger *zap.Logger, db *sql.DB, manager *tasks.Manager) http.Handler {
	// DatabaseMiddleware -> ManagerMiddleware -> LoggingMiddleware -> HandlerFunc
	return &middleware.DatabaseMiddleware{
		Database: db,
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
func setupRouter(logger *zap.Logger, db *sql.DB, manager *tasks.Manager) http.Handler {
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
		handler := setupMiddlewares(routeConfig.handler, logger, db, manager)
		router.Handle(path, handler).Methods(routeConfig.methods...)
	}

	return cors.Default().Handler(router)
}

func startHTTPServer(logger *zap.Logger, db *sql.DB, manager *tasks.Manager) {
	logger.Info("Starting the HTTP server...")
	if err := http.ListenAndServe(config.HTTPAddress, setupRouter(logger, db, manager)); err != nil {
		logger.Fatal(err.Error())
	}
}
