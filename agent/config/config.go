package config

import (
	"os"
	"strconv"
)

var BackendAddress = "localhost:5000" // gRPC адрес оркестратора

func init() {
	url := os.Getenv("BACKEND_ADDRESS")

	if url != "" {
		BackendAddress = url
	}
}

// GetComputingPower Получение COMPUTING_POWER из переменных среды
func GetComputingPower() (int, error) {
	return strconv.Atoi(os.Getenv("COMPUTING_POWER"))
}
