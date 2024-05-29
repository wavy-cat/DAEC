package storage

import "sync"

// Storage примитивная база данных типа "ключ: значение"
type Storage struct {
	database map[string]interface{}
	mutex    sync.RWMutex
}

// NewStorage Конструктор для структуры Storage
func NewStorage() *Storage {
	return &Storage{
		database: make(map[string]interface{}),
		mutex:    sync.RWMutex{},
	}
}
