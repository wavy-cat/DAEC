package main

import (
	"backend/pkg/postfix"
	"fmt"
	"math"
	"time"
)

type solver struct {
	SleepDuration time.Duration
}

func (s *solver) Addition(num1, num2 *postfix.LazyFloat) *postfix.LazyFloat {
	var float postfix.LazyFloat

	go func() {
		n1, n2 := num1.GetValue(), num2.GetValue()
		time.Sleep(s.SleepDuration)
		float.SetValue(n1 + n2)
	}()

	return &float
}

func (s *solver) Subtraction(num1, num2 *postfix.LazyFloat) *postfix.LazyFloat {
	var float postfix.LazyFloat

	go func() {
		n1, n2 := num1.GetValue(), num2.GetValue()
		time.Sleep(s.SleepDuration)
		float.SetValue(n1 - n2)
	}()

	return &float
}

func (s *solver) Division(num1, num2 *postfix.LazyFloat) *postfix.LazyFloat {
	var float postfix.LazyFloat

	go func() {
		n1, n2 := num1.GetValue(), num2.GetValue()
		time.Sleep(s.SleepDuration)
		float.SetValue(n1 / n2)
	}()

	return &float
}

func (s *solver) Multiplication(num1, num2 *postfix.LazyFloat) *postfix.LazyFloat {
	var float postfix.LazyFloat

	go func() {
		n1, n2 := num1.GetValue(), num2.GetValue()
		time.Sleep(s.SleepDuration)
		float.SetValue(n1 * n2)
	}()

	return &float
}

func (s *solver) Exponentiation(num1, num2 *postfix.LazyFloat) *postfix.LazyFloat {
	var float postfix.LazyFloat

	go func() {
		n1, n2 := num1.GetValue(), num2.GetValue()
		time.Sleep(s.SleepDuration)
		float.SetValue(math.Pow(n1, n2))
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

func calc(exp string, wait float64, sleep time.Duration) {
	start := time.Now()
	postfixNotation, err := postfix.Convertor(exp)
	if err != nil {
		fmt.Println("convertor", err)
		return
	}
	//showExp(postfixNotation)
	result := postfix.Calculate(postfixNotation, &solver{SleepDuration: sleep})

	for !result.IsDone {
		time.Sleep(100 * time.Millisecond)
	}

	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}
	duration := time.Since(start)

	var text string
	if result.Result == wait {
		text = "PASS"
	} else {
		text = "FAIL"
	}

	fmt.Println(result.Result, wait, text, duration)
}

func main() {
	// TODO: Перенести эти тестики в examples
	// ((2 + 2 * 2) ^ 2 + 4) / 2 ^ 2 - (-100 + 50 * 2) - 2
	exps := []struct {
		expression string
		expect     float64
	}{
		{"40-2", 38},
		{"-40-2", -42},
		{"40- 2", 38},
		{"40 - ( -2 )", 42},
		{"40-(-2)", 42},
		{"2 + 2 * 2", 6},
		{"10 * 5 + 10 * 2", 70},
		{"((2 + 2 * 2) ^ 2 + 4) / 2 ^ 2", 10},
		{"((2 + 2 * 2) ^ 2 + 4) / 2 ^ 2 - (-100 + 50 * 2) - 2", 8},
		{"((2 + 2 * 2) ^ 2 + 4) / 2 ^ 2 - (-100 + 50 * 2) - 2 + (2 * 2)^2", 24},
		{"5^(-7)/5^(-6)*5^3", 25},
		{"(11^(-2))/(22^(-2))", 4},
		{"-100.25 + 50 * 0.5", -75.25},
		{"192/(26-14)", 16},
		{"(176+343)-243", 276},
		{"25*(63-741/19)", 600},
		{"(900-7218/9)*12", 1176},
		{"3926/13*8+2584", 5000},
		{"690-2944/64*15", 0},
		{"(34-6-6*2)/2", 8},
		{"(25*82)^8-(99/4)", 3.1191114176253905e+26},
		{"10.5^2-2*5", 100.25},
	}
	for _, exp := range exps {
		calc(exp.expression, exp.expect, time.Second)
	}
}
