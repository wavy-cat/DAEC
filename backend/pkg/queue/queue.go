package queue

import "sync"

type Queue[T any] struct {
	mutex    sync.Mutex
	elements []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{elements: make([]T, 0)}
}

// Enqueue Добавляет новый элемент в очередь
func (q *Queue[T]) Enqueue(element T) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.elements = append(q.elements, element)
}

// Dequeue Забирает элемент из очереди
func (q *Queue[T]) Dequeue() (T, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.isEmpty() {
		var null T
		return null, false
	}
	element := q.elements[0]
	q.elements = q.elements[1:]
	return element, true
}

func (q *Queue[T]) isEmpty() bool {
	return len(q.elements) == 0
}
