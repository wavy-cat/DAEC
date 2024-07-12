package stack

import "testing"

func TestStack(t *testing.T) {
	stack := NewStack[int]()

	t.Run("Push", func(t *testing.T) {
		stack.Push(10)
		stack.Push(20)
	})

	t.Run("Peek", func(t *testing.T) {
		val, ok := stack.Peek()
		const except = 20

		if !ok {
			t.Fatal(`stack.Peek() !ok`)
		}

		if val != except {
			t.Fatalf(`stack.Peek() = %d, want %d`, val, except)
		}
	})

	t.Run("Pop", func(t *testing.T) {
		val1, ok1 := stack.Pop()
		val2, ok2 := stack.Pop()
		except1, except2 := 20, 10

		if !ok1 || !ok2 {
			t.Fatalf(`stack.Pop() !ok: %t, %t`, ok1, ok2)
		}

		if val1 != except1 {
			t.Fatalf(`stack.Pop() (1) = %d, want %d`, val1, except1)
		}
		if val2 != except2 {
			t.Fatalf(`stack.Pop() (2) = %d, want %d`, val2, except2)
		}
	})

	t.Run("Size", func(t *testing.T) {
		size := stack.Size()
		const except = 0

		if size != except {
			t.Fatalf(`stack.Size() = %d, want %d`, size, except)
		}
	})

	t.Run("Peek from empty stack", func(t *testing.T) {
		_, ok := stack.Peek()
		const except = false

		if ok != except {
			t.Fatalf(`stack.Peek() ok = %t, want %t`, ok, except)
		}
	})

	t.Run("Pop from empty stack", func(t *testing.T) {
		_, ok := stack.Pop()
		const except = false

		if ok != except {
			t.Fatalf(`stack.Pop() ok = %t, want %t`, ok, except)
		}
	})
}
