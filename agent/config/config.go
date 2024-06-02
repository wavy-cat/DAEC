package config

import (
	"os"
	"strconv"
)

const BackendUrl = "http://localhost:8080" // URL оркестратора

// GetComputingPower Получение COMPUTING_POWER из переменных среды
func GetComputingPower() (int, error) {
	return strconv.Atoi(os.Getenv("COMPUTING_POWER"))
}
