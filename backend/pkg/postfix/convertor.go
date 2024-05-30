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

func higherPrecedence(op1, op2 rune) bool {
	if op2 == '(' {
		return false
	}
	if operatorPrecedence[op1] > operatorPrecedence[op2] {
		return true
	}
	if operatorPrecedence[op1] == operatorPrecedence[op2] &&
		operatorAssociativity[op1] == "left" {
		return true
	}
	return false
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
		switch token.(type) {
		case float64:
			outputQueue = append(outputQueue, token)
		default:
			switch token {
			case '+', '-', '/', '*', '^':
				if op2, ok := stack.Peek(); ok && higherPrecedence(token.(rune), op2) {
					stack.Pop()
					outputQueue = append(outputQueue, op2)
				}
				stack.Push(token.(rune))
			case '(':
				stack.Push(token.(rune))
			case ')':
				for op, ok := stack.Peek(); op != '('; op, ok = stack.Peek() {
					if !ok {
						return nil, errors.New("the stack ran out too early. error in expression")
					}
					stack.Pop()
					outputQueue = append(outputQueue, op)
				}
				stack.Pop() // выкидываем открывающуюся скобку из стека
			}
		}
	}

	if val, ok := stack.Peek(); ok {
		if val == '(' {
			return nil, errors.New("missing parenthesis in expression")
		}
		for token, ok := stack.Pop(); ok; token, ok = stack.Pop() {
			outputQueue = append(outputQueue, token)
		}
	}

	return outputQueue, nil
}
