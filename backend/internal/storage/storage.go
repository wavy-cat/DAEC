package storage

import "sync"

// Storage примитивная база данных типа "ключ: значение"
type Storage[T any] struct {
	database map[any]T
	mutex    sync.RWMutex
}

// NewStorage Конструктор для структуры Storage
func NewStorage[T any]() *Storage[T] {
	return &Storage[T]{
		database: make(map[any]T),
		mutex:    sync.RWMutex{},
	}
}

// Set Добавление значения по ключу. Перезаписывает существующее значение при том же ключе.
func (s *Storage[T]) Set(key any, value T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.database[key] = value
}

// Get Получение значения по ключу.
func (s *Storage[T]) Get(key any) (T, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	val, ok := s.database[key]
	return val, ok
}

// GetAll Получение всех записей в виде слайса структуры Key any, Value T.
func (s *Storage[T]) GetAll() []T {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	r := make([]T, 0)

	for _, val := range s.database {
		r = append(r, val)
	}

	return r
}

// Del Удаление значения по ключу.
func (s *Storage[T]) Del(key any) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.database, key)
}
