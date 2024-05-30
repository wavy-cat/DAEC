package main

import (
	"backend/pkg/postfix"
	"fmt"
	"math"
	"time"
)

type solver struct {
}

func (s *solver) Addition(num1, num2 *postfix.LazyFloat) *postfix.LazyFloat {
	var float postfix.LazyFloat

	go func() {
		float.SetValue(num1.GetValue() + num2.GetValue())
	}()

	return &float
}

func (s *solver) Subtraction(num1, num2 *postfix.LazyFloat) *postfix.LazyFloat {
	var float postfix.LazyFloat

	go func() {
		float.SetValue(num1.GetValue() - num2.GetValue())
	}()

	return &float
}

func (s *solver) Division(num1, num2 *postfix.LazyFloat) *postfix.LazyFloat {
	var float postfix.LazyFloat

	go func() {
		float.SetValue(num1.GetValue() / num2.GetValue())
	}()

	return &float
}

func (s *solver) Multiplication(num1, num2 *postfix.LazyFloat) *postfix.LazyFloat {
	var float postfix.LazyFloat

	go func() {
		float.SetValue(num1.GetValue() * num2.GetValue())
	}()

	return &float
}

func (s *solver) Exponentiation(num1, num2 *postfix.LazyFloat) *postfix.LazyFloat {
	var float postfix.LazyFloat

	go func() {
		float.SetValue(math.Pow(num1.GetValue(), num2.GetValue()))
	}()

	return &float
}

func showExp(exp []any) {
	for _, token := range exp {
		switch token.(type) {
		case float64:
			fmt.Print(token)
		case rune:
			fmt.Print(string(token.(rune)))
		}
		fmt.Print(" ")
	}
	fmt.Println()
}

func calc(exp string, wait float64) {
	postfixNotation, err := postfix.Convertor(exp)
	if err != nil {
		fmt.Println("convertor", err)
		return
	}
	//showExp(postfixNotation)
	result := postfix.Calculate(postfixNotation, &solver{})

	for !result.IsDone {
		time.Sleep(100 * time.Millisecond)
	}

	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}

	var text string
	if result.Result == wait {
		text = "PASS"
	} else {
		text = "FAIL"
	}

	fmt.Println(result.Result, wait, text)
}

func main() {
	// TODO: Перенести эти тестики в examples
	// ((2 + 2 * 2) ^ 2 + 4) / 2 ^ 2 - (-100 + 50 * 2) - 2
	calc("((2 + 2 * 2)) ", 6)
	calc("((2 + 2 * 2) ^ 2 + 4) / 2 ^ 2", 10)
	calc("((2 + 2 * 2) ^ 2 + 4) / 2 ^ 2 - (-100 + 50 * 2) - 2", 8)
	calc("((2 + 2 * 2) ^ 2 + 4) / 2 ^ 2 - (-100 + 50 * 2) - 2 + (2 * 2)^2", 24)
	calc("5^(-7)/5^(-6)*5^3", 25)
	calc("(11^(-2))/(22^(-2))", 4)
	calc("-100.25 + 50 * 0.5", -75.25)
	calc("192/(26-14)", 16)
	calc("(176+343)- 243", 276)
	calc("25*(63-741/19)", 600)
	calc("(900-7218/9)*12", 1176)
	calc("3926/13*8+2584", 5000)
	calc("690-2944/64*15", 0)
	calc("(34-6-6*2)/2", 8)
	calc("(25*82)^8-(99/4)", 3.1191114176253905e+26)
}
