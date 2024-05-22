package middleware

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// LoggingMiddleware - Middleware для передачи объекта логгера в обработчик
type LoggingMiddleware struct {
	Logger *zap.Logger  // Указатель на объект логгера
	Next   http.Handler // Функция, вызываемая middleware, которая будет обрабатывать http запрос
}

func (mw *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	r = r.WithContext(context.WithValue(r.Context(), "logger", mw.Logger))
	mw.Next.ServeHTTP(w, r)

	duration := time.Since(start)
	mw.Logger.Info("HTTP request",
		zap.String("path", r.URL.Path),
		zap.Duration("duration", duration),
	)
}
