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
					// Проверка то что минус не стоит в конце выражения
					return nil, errors.New("operator detected at end of expression")
				}

				// Если минус находится в начале выражения или перед ним не стоит цифра, но anyway идёт после,
				// то это унарный минус. Если нет, то обычный минус.
				if !((i == 0 || (!unicode.IsDigit(rune(expression[i-1])) && rune(expression[i-1]) != ')')) && unicode.IsDigit(rune(expression[i+1]))) {
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

			n, _ := strconv.ParseFloat(num, 64) // Предположим что нет такой ситуации чтобы была ошибка.
			// А если найдёте — киньте Issue, будет интересно посмотреть.

			if negative {
				n = -n
			}

			tokens = append(tokens, n)
		}
	}

	return tokens, nil
}
