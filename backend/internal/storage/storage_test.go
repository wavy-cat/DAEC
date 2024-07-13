package storage

import (
	"reflect"
	"testing"
)

func Test(t *testing.T) {
	storage := NewStorage[string]()

	t.Run("Set value", func(t *testing.T) {
		storage.Set("foo", "bar")
	})

	t.Run("Get value", func(t *testing.T) {
		val, ok := storage.Get("foo")
		const except = "bar"

		if !ok {
			t.Fatal(`storage.Get("foo") !ok`)
		}

		if val != except {
			t.Fatalf(`storage.Get("foo") = %s, want %s`, val, except)
		}
	})

	t.Run("Get unknown value", func(t *testing.T) {
		_, ok := storage.Get("hello")

		if ok {
			t.Fatal(`storage.Get("foo") ok`)
		}
	})

	t.Run("Get values", func(t *testing.T) {
		val := storage.GetAll()
		except := []string{"bar"}

		if !reflect.DeepEqual(val, except) {
			t.Fatalf(`storage.Get("foo") = %q, want %q`, val, except)
		}
	})

	t.Run("Del value", func(t *testing.T) {
		storage.Del("foo")
		_, ok := storage.Get("foo")

		if ok {
			t.Fatal(`storage.Get("foo") !ok after storage.Del("foo")`)
		}
	})
}
