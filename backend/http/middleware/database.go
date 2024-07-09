package middleware

import (
	"context"
	"database/sql"
	"net/http"
)

// DatabaseMiddleware - Middleware для передачи объекта базы данных (Storage) в обработчик
type DatabaseMiddleware struct {
	Database *sql.DB      // Указатель на объект базы данных
	Next     http.Handler // Функция, вызываемая middleware, которая будет обрабатывать http запрос
}

func (mw *DatabaseMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(r.Context(), "database", mw.Database))
	mw.Next.ServeHTTP(w, r)
}
