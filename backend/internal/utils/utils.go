package utils

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
