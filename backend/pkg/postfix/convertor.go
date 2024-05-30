package postfix

import (
	stk "backend/pkg/stack"
	"errors"
)

var operatorPrecedence = map[rune]int{
	'+': 1,
	'-': 1,
	'/': 2,
	'*': 2,
	'^': 3,
}

var operatorAssociativity = map[rune]string{
	'+': "left",
	'-': "left",
	'/': "left",
	'*': "left",
	'^': "right",
}

// Convertor Переводит выражение в инфиксной нотации в слайс токенов постфиксной записи.
// Возвращаемый слайс состоит из float64 и rune.
func Convertor(expression string) ([]any, error) {
	stack := stk.NewStack[rune]()
	outputQueue := make([]any, 0)
	tokens, err := SplitExpressionToTokens(expression)

	if err != nil {
		return nil, err
	}

	for _, token := range tokens {
		switch token := token.(type) {
		case float64:
			outputQueue = append(outputQueue, token)
		case rune:
			switch token {
			case '+', '-', '/', '*', '^':
				for {
					op2, ok := stack.Peek()
					if !ok || op2 == '(' {
						break
					}
					if (operatorAssociativity[token] == "left" && operatorPrecedence[token] <= operatorPrecedence[op2]) ||
						(operatorAssociativity[token] == "right" && operatorPrecedence[token] < operatorPrecedence[op2]) {
						stack.Pop()
						outputQueue = append(outputQueue, op2)
					} else {
						break
					}
				}
				stack.Push(token)
			case '(':
				stack.Push(token)
			case ')':
				for {
					op, ok := stack.Peek()
					if !ok {
						return nil, errors.New("the stack ran out too early. error in expression")
					}
					if op == '(' {
						break
					}
					stack.Pop()
					outputQueue = append(outputQueue, op)
				}
				stack.Pop() // выкидываем открывающуюся скобку из стека
			}
		}
	}

	// После завершения цикла проверяем оставшиеся операторы в стеке
	for {
		token, ok := stack.Pop()
		if !ok {
			break
		}
		if token == '(' {
			return nil, errors.New("missing parenthesis in expression")
		}
		outputQueue = append(outputQueue, token)
	}

	return outputQueue, nil
}
