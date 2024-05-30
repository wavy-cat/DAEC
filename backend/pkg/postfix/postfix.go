package postfix

import (
	stk "backend/pkg/stack"
	"errors"
	"time"
)

// LazyFloat - "ленивый" или отложенный float64.
// Значение может быть пустое (ещё вычисляться).
// О готовности можно узнать из IsDone.
type LazyFloat struct {
	IsDone bool    // Статус значения (посчиталось ещё или нет)
	Value  float64 // Значение
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

	go func(result *CalcResult) {
		stack := stk.NewStack[*LazyFloat]()

		for _, token := range postfixNotation {
			switch token.(type) {
			case float64:
				stack.Push(&LazyFloat{Value: token.(float64), IsDone: true})
			case rune:
				num1, ok1 := stack.Pop()
				num2, ok2 := stack.Pop()
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

		num, ok := stack.Pop()
		if ok {
			result.Error = errors.New("error getting value from stack. internal server error")
			result.IsDone = true
			return
		}

		for !num.IsDone {
			time.Sleep(100 * time.Millisecond)
		}

		result.Result = num.Value
		result.IsDone = true
	}(&result)

	return &result
}
