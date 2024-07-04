package middleware

import (
	"context"
	"github.com/wavy-cat/DAEC/backend/internal/tasks"
	"net/http"
)

// ManagerMiddleware - Middleware для передачи объекта менеджера задач (tasks.Manager) в обработчик.
type ManagerMiddleware struct {
	Manager *tasks.Manager // Указатель на tasks.Manager
	Next    http.Handler   // Функция, вызываемая middleware, которая будет обрабатывать http запрос
}

func (mw *ManagerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(r.Context(), "manager", mw.Manager))
	mw.Next.ServeHTTP(w, r)
}
