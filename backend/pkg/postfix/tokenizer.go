package postfix

import (
	"errors"
	"strconv"
	"unicode"
)

// SplitExpressionToTokens Разбивает выражение в инфиксной нотации на токены.
// Возвращает слайс из float64 и rune.
func SplitExpressionToTokens(expression string) ([]any, error) {
	tokens := make([]any, 0, 1)

	for i := 0; i < len(expression); i++ {
		char := rune(expression[i])

		switch char {
		case '(', ')':
			tokens = append(tokens, char)
		case '+', '/', '*', '^':
			if len(expression)-1 == i {
				return nil, errors.New("operator detected at end of expression")
			}
			tokens = append(tokens, char)
		default:
			var negative bool

			if char == '-' {
				if len(expression)-1 == i {
					return nil, errors.New("operator detected at end of expression")
				}
				if !unicode.IsDigit(rune(expression[i+1])) {
					tokens = append(tokens, char)
					continue
				}

				negative = true
				i++
				char = rune(expression[i])
			} else {
				if !unicode.IsDigit(rune(expression[i])) {
					continue
				}
			}

			var num string

			for unicode.IsDigit(char) || char == '.' {
				num += string(char)
				i++
				if len(expression) == i {
					break
				}
				char = rune(expression[i])
			}
			i--

			n, err := strconv.ParseFloat(num, 64)
			if err != nil {
				return nil, err
			}

			if negative {
				n = -n
			}

			tokens = append(tokens, n)
		}
	}

	return tokens, nil
}
