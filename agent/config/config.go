package config

import (
	"os"
	"strconv"
)

var BackendUrl = "http://localhost" // URL оркестратора

func init() {
	url := os.Getenv("BACKEND_URL")

	if url != "" {
		BackendUrl = url
	}
}

// GetComputingPower Получение COMPUTING_POWER из переменных среды
func GetComputingPower() (int, error) {
	return strconv.Atoi(os.Getenv("COMPUTING_POWER"))
}
