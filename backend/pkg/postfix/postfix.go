package postfix

import (
	"errors"
	stk "github.com/wavy-cat/DAEC/backend/pkg/stack"
	"time"
)

// LazyFloat - "ленивый" или отложенный float64.
// Значение может быть пустое (ещё вычисляться).
type LazyFloat struct {
	isDone bool    // Статус значения (посчиталось ещё или нет)
	IsFail bool    // True, если значение не удалось вычислить
	value  float64 // Значение
}

// GetValue отдаёт значение, когда оно готово.
func (f *LazyFloat) GetValue() float64 {
	for !f.isDone {
		time.Sleep(10 * time.Millisecond)
	}
	return f.value
}

// SetValue устанавливает значение и его готовность.
func (f *LazyFloat) SetValue(num float64) {
	f.value = num
	f.isDone = true
}

func (f *LazyFloat) SetFail() {
	f.IsFail = true
	f.isDone = true
}

// MathSolver Интерфейс, который реализует базовые арифметические операции.
// Методы должны возвращать LazyFloat незамедлительно, поэтому значение должно вычисляться в фоне.
type MathSolver interface {
	Addition(num1, num2 *LazyFloat) *LazyFloat       // Сложение (+)
	Subtraction(num1, num2 *LazyFloat) *LazyFloat    // Вычитание (-)
	Division(num1, num2 *LazyFloat) *LazyFloat       // Деление (/)
	Multiplication(num1, num2 *LazyFloat) *LazyFloat // Умножение (*)
	Exponentiation(num1, num2 *LazyFloat) *LazyFloat // Возведение в степень (^)
}

// CalcResult Структура, содержащая результат вычисления выражения.
type CalcResult struct {
	Result float64 // Результат
	IsDone bool    // Статус вычисления (вычислено ли)
	Error  error   // Ошибка, если есть
}

// Calculate Вычисляет значение выражения в постфиксной записи.
// Для работы требуется слайс токенов постфиксной нотации (можно получить из Convertor).
// А также структура, реализующая интерфейс MathSolver.
func Calculate(postfixNotation []any, solver MathSolver) *CalcResult {
	var result CalcResult

	go func(result *CalcResult, solver MathSolver) {
		stack := stk.NewStack[*LazyFloat]()

		for _, token := range postfixNotation {
			switch token.(type) {
			case float64:
				stack.Push(&LazyFloat{value: token.(float64), isDone: true})
			case rune:
				num2, ok1 := stack.Pop()
				num1, ok2 := stack.Pop()
				if !ok1 || !ok2 {
					result.Error = errors.New("error getting value from stack. " +
						"perhaps there is an error in the expression")
					result.IsDone = true
					return
				}

				switch token {
				case '+':
					stack.Push(solver.Addition(num1, num2))
				case '-':
					stack.Push(solver.Subtraction(num1, num2))
				case '/':
					stack.Push(solver.Division(num1, num2))
				case '*':
					stack.Push(solver.Multiplication(num1, num2))
				case '^':
					stack.Push(solver.Exponentiation(num1, num2))
				}
			}
		}

		if stack.Size() != 1 {
			result.Error = errors.New("incorrect number of elements in the stack. " +
				"perhaps there is an error in the expression")
			result.IsDone = true
			return
		}

		num, _ := stack.Pop() // Тут должно быть всё ок, т.к. раньше этот момент проверили.

		result.Result = num.GetValue()
		if num.IsFail {
			result.Error = errors.New("an error occurred during calculation")
		}
		result.IsDone = true
	}(&result, solver)

	return &result
}
