package storage

import "sync"

// Storage примитивная база данных типа "ключ: значение"
type Storage struct {
	database map[string]any
	mutex    sync.RWMutex
}

// NewStorage Конструктор для структуры Storage
func NewStorage() *Storage {
	return &Storage{
		database: make(map[string]any),
		mutex:    sync.RWMutex{},
	}
}

func (s *Storage) Set(key string, value any) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.database[key] = value
}

func (s *Storage) Get(key string) any {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.database[key]
}

func (s *Storage) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.database, key)
}
