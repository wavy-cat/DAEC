package postfix

import (
	"fmt"
	"math"
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

	t.Run("Operator at end", func(t *testing.T) {
		t.Parallel()

		_, err := Convertor("2 + 2 -")

		if err == nil {
			t.Fatalf(`Convertor("2 + 2 -") = nil, want error`)
		}
	})
}

func TestCalculate(t *testing.T) {
	t.Run("Plan Success", func(t *testing.T) {
		t.Parallel()

		notation, _ := Convertor("2 + 2 * 2 / 2 ^ 2 - 2")
		result := Calculate(notation, &solver{})
		except := 1.0

		for !result.IsDone {
		}

		if result.Error != nil {
			t.Fatalf(`Calculate(notation, &solver{}) has error %s, want %f`, result.Error.Error(), except)
		}

		if result.Result != except {
			t.Fatalf(`Calculate(notation, &solver{}) = %f, want %f`, result.Result, except)
		}
	})

	t.Run("Error getting value from stack", func(t *testing.T) {
		t.Parallel()

		notation, _ := Convertor("2 + + 2")
		result := Calculate(notation, &solver{})

		for !result.IsDone {
		}

		if result.Error == nil {
			t.Fatalf(`Calculate(notation, &solver{}).Error = nil, want error`)
		}
	})

	t.Run("Incorrect parentheses", func(t *testing.T) {
		t.Parallel()

		notation, _ := Convertor("( 2 + 3")
		result := Calculate(notation, &solver{})

		for !result.IsDone {
		}

		if result.Error == nil {
			t.Fatalf(`Calculate(notation, &solver{}).Error = nil, want error`)
		}
	})

	t.Run("Set fail in solver", func(t *testing.T) {
		t.Parallel()

		notation, _ := Convertor("2 + 3")
		result := Calculate(notation, &solver{NeedFailed: true})

		for !result.IsDone {
		}

		if result.Error == nil {
			t.Fatalf(`Calculate(notation, &solver{}).Error = nil, want error`)
		}
	})

	t.Run("Problem from stack (line 100)", func(t *testing.T) {
		t.Parallel()

		notation, _ := Convertor("(+)")
		result := Calculate(notation, &solver{})

		for !result.IsDone {
		}

		fmt.Println(result.Error)

		if result.Error == nil {
			t.Fatalf(`Calculate(notation, &solver{}).Error = nil, want error`)
		}
	})
}

// Реализация интерфейса MathSolver
type solver struct {
	NeedFailed bool // Нужен ли Fail
}

func (s *solver) Addition(num1, num2 *LazyFloat) *LazyFloat {
	var lf LazyFloat
	go func() {
		switch s.NeedFailed {
		case false:
			n1, n2 := num1.GetValue(), num2.GetValue()
			if num1.IsFail || num2.IsFail {
				lf.SetFail()
				return
			}
			lf.SetValue(n1 + n2)
		case true:
			lf.SetFail()
		}
	}()
	return &lf
}

func (s *solver) Subtraction(num1, num2 *LazyFloat) *LazyFloat {
	var lf LazyFloat
	go func() {
		switch s.NeedFailed {
		case false:
			n1, n2 := num1.GetValue(), num2.GetValue()
			if num1.IsFail || num2.IsFail {
				lf.SetFail()
				return
			}
			lf.SetValue(n1 - n2)
		case true:
			lf.SetFail()
		}
	}()
	return &lf
}

func (s *solver) Division(num1, num2 *LazyFloat) *LazyFloat {
	var lf LazyFloat
	go func() {
		switch s.NeedFailed {
		case false:
			n1, n2 := num1.GetValue(), num2.GetValue()
			if num1.IsFail || num2.IsFail {
				lf.SetFail()
				return
			}
			lf.SetValue(n1 / n2)
		case true:
			lf.SetFail()
		}
	}()
	return &lf
}

func (s *solver) Multiplication(num1, num2 *LazyFloat) *LazyFloat {
	var lf LazyFloat
	go func() {
		switch s.NeedFailed {
		case false:
			n1, n2 := num1.GetValue(), num2.GetValue()
			if num1.IsFail || num2.IsFail {
				lf.SetFail()
				return
			}
			lf.SetValue(n1 * n2)
		case true:
			lf.SetFail()
		}
	}()
	return &lf
}

func (s *solver) Exponentiation(num1, num2 *LazyFloat) *LazyFloat {
	var lf LazyFloat
	go func() {
		switch s.NeedFailed {
		case false:
			n1, n2 := num1.GetValue(), num2.GetValue()
			if num1.IsFail || num2.IsFail {
				lf.SetFail()
				return
			}
			lf.SetValue(math.Pow(n1, n2))
		case true:
			lf.SetFail()
		}
	}()
	return &lf
}
