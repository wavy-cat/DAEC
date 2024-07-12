package queue

import (
	"testing"
)

func TestQueue(t *testing.T) {
	queue := NewQueue[int]()

	t.Run("Enqueue", func(t *testing.T) {
		queue.Enqueue(10)
		queue.Enqueue(20)
	})

	t.Run("Dequeue", func(t *testing.T) {
		r1, ok1 := queue.Dequeue()
		r2, ok2 := queue.Dequeue()
		exceptR1, exceptR2 := 10, 20

		if !ok1 || !ok2 {
			t.Fatalf(`Failed to get value in queue.Dequeue(): %t, %t`, ok1, ok2)
		}

		if r1 != exceptR1 {
			t.Fatalf(`queue.Dequeue() (r1) = %d, want %d`, r1, exceptR1)
		}
		if r2 != exceptR2 {
			t.Fatalf(`queue.Dequeue() (r2) = %d, want %d`, r2, exceptR2)
		}
	})

	t.Run("Queue size", func(t *testing.T) {
		isEmpty := queue.isEmpty()
		const except = true

		if isEmpty != except {
			t.Fatalf(`queue.isEmpty() = %t, want %t`, isEmpty, except)
		}
	})

	t.Run("Dequeue from empty queue", func(t *testing.T) {
		_, ok := queue.Dequeue()
		const except = false

		if ok != except {
			t.Fatalf(`queue.Dequeue() = %t, want %t`, ok, except)
		}
	})
}
