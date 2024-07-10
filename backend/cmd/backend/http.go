package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/wavy-cat/DAEC/backend/http/handler/calculate"
	"github.com/wavy-cat/DAEC/backend/http/handler/expressions"
	"github.com/wavy-cat/DAEC/backend/http/handler/user"
	"github.com/wavy-cat/DAEC/backend/http/middleware"
	"github.com/wavy-cat/DAEC/backend/internal/config"
	"github.com/wavy-cat/DAEC/backend/internal/tasks"
	"go.uber.org/zap"
	"net/http"
)

// setupBaseMiddlewares оборачивает http.HandlerFunc в базовые мидлвари (Database и Logging).
func setupBaseMiddlewares(handlerFunc http.HandlerFunc, logger *zap.Logger, db *sql.DB) http.Handler {
	// DatabaseMiddleware -> LoggingMiddleware -> HandlerFunc
	return &middleware.DatabaseMiddleware{
		Database: db,
		Next: &middleware.LoggingMiddleware{
			Logger: logger,
			Next:   handlerFunc,
		},
	}
}

// setupFullMiddlewares оборачивает http.HandlerFunc во все мидлвари.
func setupFullMiddlewares(handlerFunc http.HandlerFunc, logger *zap.Logger, db *sql.DB, manager *tasks.Manager) http.Handler {
	// Auth → Database → Manager → Logging → HandlerFunc
	return &middleware.AuthMiddleware{
		Database: db,
		Logger:   logger,
		Next: &middleware.DatabaseMiddleware{
			Database: db,
			Next: &middleware.ManagerMiddleware{
				Manager: manager,
				Next: &middleware.LoggingMiddleware{
					Logger: logger,
					Next:   handlerFunc,
				},
			},
		},
	}
}

// setupRouter создаёт новый экземпляр mux.Router и настраивает конфигурацию маршрутизации.
func setupRouter(logger *zap.Logger, db *sql.DB, manager *tasks.Manager) http.Handler {
	router := mux.NewRouter()

	// Здесь настраивается маршрутизация aka указание эндпойнтов сервера
	routes := map[string]struct {
		handler     http.HandlerFunc
		method      string
		middlewares string
	}{
		// Указываются тут, если что.
		// Ключ — путь, значение — структура из объекта HandlerFunc и разрешённых методов.
		"/api/v1/calculate": {
			calculate.Handler,
			"POST",
			"full",
		},
		"/api/v1/expressions": {
			expressions.Handler,
			"GET",
			"full",
		},
		"/api/v1/expressions/{id}": {
			expressions.HandlerById,
			"GET",
			"full",
		},
		"/api/v1/register": {
			user.RegisterHandler,
			"POST",
			"base",
		},
		"/api/v1/login": {
			user.LoginHandler,
			"POST",
			"base",
		},
	}

	for path, routeConfig := range routes {
		// Чтобы не оборачивать функции в мидлвари вручную используется setupMiddlewares
		var handler http.Handler
		switch routeConfig.middlewares {
		case "full":
			handler = setupFullMiddlewares(routeConfig.handler, logger, db, manager)
		default:
			handler = setupBaseMiddlewares(routeConfig.handler, logger, db)
		}

		router.Handle(path, handler).Methods(routeConfig.method)
	}

	return cors.Default().Handler(router)
}

func startHTTPServer(logger *zap.Logger, db *sql.DB, manager *tasks.Manager) {
	logger.Info("Starting the HTTP server...")
	if err := http.ListenAndServe(config.HTTPAddress, setupRouter(logger, db, manager)); err != nil {
		logger.Fatal(err.Error())
	}
}
