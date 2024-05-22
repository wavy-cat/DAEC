package middleware

import (
	"context"
	"net/http"
	"sync"
)

// DatabaseMiddleware - Middleware для передачи объекта базы данных в обработчик
type DatabaseMiddleware struct {
	Database *map[string]interface{} // Указатель на мапу, которая будет выступать в роли базы данных
	Lock     *sync.RWMutex           // Mutex для синхронизации
	Next     http.Handler            // Функция, вызываемая middleware, которая будет обрабатывать http запрос
}

func (mw *DatabaseMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(r.Context(), "database", mw.Database))
	r = r.WithContext(context.WithValue(r.Context(), "lock", &mw.Lock))
	mw.Next.ServeHTTP(w, r)
}
