package config

import (
	"os"
	"strconv"
)

const ServerAddress = ":80"

var (
	TimeSubtractionMs     int
	TimeMultiplicationsMs int
	TimeDivisionsMs       int
	TimeExponentiationMs  int
	TimeAdditionMs        int
)

func init() {
	var err error

	TimeAdditionMs, err = strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
	if err != nil {
		panic(err)
	}

	TimeSubtractionMs, err = strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
	if err != nil {
		panic(err)
	}

	TimeMultiplicationsMs, err = strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
	if err != nil {
		panic(err)
	}

	TimeDivisionsMs, err = strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
	if err != nil {
		panic(err)
	}

	TimeExponentiationMs, err = strconv.Atoi(os.Getenv("TIME_EXPONENTIATION_MS"))
	if err != nil {
		panic(err)
	}
}
