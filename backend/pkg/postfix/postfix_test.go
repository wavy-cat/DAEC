package postfix

import (
	"reflect"
	"testing"
)

func TestTokenizer(t *testing.T) {
	t.Run("Short", func(t *testing.T) {
		t.Parallel()

		result, err := SplitExpressionToTokens("2 + 2")
		except := []any{2.0, '+', 2.0}

		if err != nil {
			t.Fatalf(`Unexpected error in SplitExpressionToTokens("2 + 2"): %v`, except)
		}
		if !reflect.DeepEqual(result, except) {
			t.Fatalf(`SplitExpressionToTokens("2 + 2") = %q, want %q`, result, except)
		}
	})

	t.Run("Long", func(t *testing.T) {
		t.Parallel()

		result, err := SplitExpressionToTokens("2 + 2 * (2 / 2) ^ 4 - 10")
		except := []any{2.0, '+', 2.0, '*', '(', 2.0, '/', 2.0, ')', '^', 4.0, '-', 10.0}

		if err != nil {
			t.Fatalf(`Unexpected error in SplitExpressionToTokens("2 + 2 * (2 / 2) ^ 4 - 10"): %v`, except)
		}
		if !reflect.DeepEqual(result, except) {
			t.Fatalf(`SplitExpressionToTokens("2 + 2 * (2 / 2) ^ 4 - 10") = %q, want %q`, result, except)
		}
	})

	t.Run("Operator at end 1", func(t *testing.T) {
		t.Parallel()

		_, err := SplitExpressionToTokens("2 + 2 -")

		if err == nil {
			t.Fatalf(`SplitExpressionToTokens("2 + 2 -") = nil, want error`)
		}
	})

	t.Run("Operator at end 2", func(t *testing.T) {
		t.Parallel()

		_, err := SplitExpressionToTokens("2 + 2 *")

		if err == nil {
			t.Fatalf(`SplitExpressionToTokens("2 + 2 *") = nil, want error`)
		}
	})

	t.Run("Negative", func(t *testing.T) {
		t.Parallel()

		result, err := SplitExpressionToTokens("-2 + 2")
		except := []any{-2.0, '+', 2.0}

		if err != nil {
			t.Fatalf(`Unexpected error in SplitExpressionToTokens("-2 + 2"): %v`, except)
		}
		if !reflect.DeepEqual(result, except) {
			t.Fatalf(`SplitExpressionToTokens("-2 + 2") = %q, want %q`, result, except)
		}
	})

	t.Run("Some else", func(t *testing.T) {
		t.Parallel()

		result, err := SplitExpressionToTokens("-2.5 - 2.5")
		except := []any{-2.5, '-', 2.5}

		if err != nil {
			t.Fatalf(`Unexpected error in SplitExpressionToTokens("-2.5 - 2"): %v`, except)
		}
		if !reflect.DeepEqual(result, except) {
			t.Fatalf(`SplitExpressionToTokens("-2.5 - 2") = %q, want %q`, result, except)
		}
	})
}

func TestConvertor(t *testing.T) {
	t.Run("Short", func(t *testing.T) {
		t.Parallel()

		result, err := Convertor("2 + 2")
		except := []any{2.0, 2.0, '+'}

		if err != nil {
			t.Fatalf(`Unexpected error in Convertor("2 + 2"): %v`, except)
		}
		if !reflect.DeepEqual(result, except) {
			t.Fatalf(`Convertor("2 + 2") = %q, want %q`, result, except)
		}
	})

	t.Run("Long", func(t *testing.T) {
		t.Parallel()

		result, err := Convertor("( 20-18 ) ^ 2 + 10 * 4")
		except := []any{20.0, 18.0, '-', 2.0, '^', 10.0, 4.0, '*', '+'}

		if err != nil {
			t.Fatalf(`Unexpected error in Convertor("2 + 2 * (2 / 2) ^ 4 - 10"): %v`, except)
		}
		if !reflect.DeepEqual(result, except) {
			t.Fatalf(`Convertor("2 + 2 * (2 / 2) ^ 4 - 10") = %q, want %q`, result, except)
		}
	})

	t.Run("Stack ran out too early", func(t *testing.T) {
		t.Parallel()

		_, err := Convertor("( 20-18 ) )")

		if err == nil {
			t.Fatalf(`Convertor("( 20-18 ) )") = nil, want error`)
		}
	})

	t.Run("Missing parenthesis in expression", func(t *testing.T) {
		t.Parallel()

		_, err := Convertor("( 20-18 ) + (")

		if err == nil {
			t.Fatalf(`Convertor("( 20-18 ) + (") = nil, want error`)
		}
	})
}

func TestCalculate(t *testing.T) {

}
