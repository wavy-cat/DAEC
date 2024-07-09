package config

import (
	"os"
	"strconv"
)

var HTTPAddress = ":80"   // Адрес HTTP сервера
var GRPCAddress = ":5000" // Адрес gRPC сервера

const DatabasePath = "database.db" // Путь до файла базы данных

var (
	TimeSubtractionMs     int
	TimeMultiplicationsMs int
	TimeDivisionsMs       int
	TimeExponentiationMs  int
	TimeAdditionMs        int
)

func setAddresses() {
	if httpAddr := os.Getenv("HTTP_ADDRESS"); httpAddr != "" {
		HTTPAddress = httpAddr
	}

	if grpcAddr := os.Getenv("GRPC_ADDRESS"); grpcAddr != "" {
		GRPCAddress = grpcAddr
	}
}

func setSleepTime() {
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

func init() {
	setSleepTime()
}
