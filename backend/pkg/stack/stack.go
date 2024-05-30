package stack

import "sync"

// Stack Простая реализация стека
type Stack[T any] struct {
	lock  sync.RWMutex
	items []interface{}
}

// NewStack Создаёт новый экземпляр структуры Stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Push добавляет элемент в стек
func (stack *Stack[T]) Push(item T) {
	stack.lock.Lock()
	defer stack.lock.Unlock()

	stack.items = append(stack.items, item)
}

// Pop удаляет элемент из вершины стека
func (stack *Stack[T]) Pop() (T, bool) {
	stack.lock.Lock()
	defer stack.lock.Unlock()

	if len(stack.items) == 0 {
		var null T
		return null, false
	}

	item := stack.items[len(stack.items)-1]
	stack.items = stack.items[0 : len(stack.items)-1]

	return item, true
}

// Size возвращает текущий размер стека
func (stack *Stack[T]) Size() int {
	stack.lock.RLock()
	defer stack.lock.RUnlock()

	return len(stack.items)
}

// Peek возвращает элемент из вершины стека, не удаляя его
func (stack *Stack[T]) Peek() (T, bool) {
	stack.lock.RLock()
	defer stack.lock.RUnlock()

	if len(stack.items) == 0 {
		var null T
		return null, false
	}

	return stack.items[len(stack.items)-1], true
}
