package middleware

import (
	"backend/internal/storage"
	"context"
	"net/http"
)

// DatabaseMiddleware - Middleware для передачи объекта базы данных (Storage) в обработчик
type DatabaseMiddleware struct {
	Storage *storage.Storage // Указатель на структуру Storage (БД)
	Next    http.Handler     // Функция, вызываемая middleware, которая будет обрабатывать http запрос
}

func (mw *DatabaseMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(r.Context(), "storage", mw.Storage))
	mw.Next.ServeHTTP(w, r)
}
