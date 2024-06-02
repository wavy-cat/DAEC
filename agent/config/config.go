package config

import (
	"os"
	"strconv"
)

const BackendUrl = "http://localhost" // URL оркестратора

// GetComputingPower Получение COMPUTING_POWER из переменных среды
func GetComputingPower() (int, error) {
	return strconv.Atoi(os.Getenv("COMPUTING_POWER"))
}
