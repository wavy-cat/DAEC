package utils

import pb "github.com/wavy-cat/DAEC/backend/proto"

// CheckCharsInString Проверяет все ли символы из строки присутствуют в слайсе chars.
// При false возвращает первую руну, которая не состоит в слайсе.
func CheckCharsInString(str string, chars []rune) (bool, rune) {
	for _, char := range str {
		contains := false
		for _, c := range chars {
			if char == c {
				contains = true
				break
			}
		}
		if !contains {
			return false, char
		}
	}
	return true, 0
}

// ParseOperation Преобразует строковую операцию в соответствующую Operation константу.
func ParseOperation(operation byte) pb.Operation {
	switch operation {
	case '+':
		return pb.Operation_ADDITION
	case '-':
		return pb.Operation_SUBTRACTION
	case '*':
		return pb.Operation_MULTIPLICATION
	case '/':
		return pb.Operation_DIVISION
	case '^':
		return pb.Operation_EXPONENTIATION
	}
	panic("unknown operation")
}

// Почему программисты не любят находиться на природе?
// Там слишком много bugов.
