package main

import (
	"agent/work"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strconv"
)

func main() {
	// Создание и запуск логгера
	logger, err := zap.NewDevelopment() // Заменить в конце на NewProduction

	if err != nil {
		fmt.Println("error initializing logger:", err)
		return
	}

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Println("error synchronizing logger:", err)
		}
	}(logger)

	logger.Info("Agent is starting")

	// Запуск worker pool (горутин, которые будут выполнять арифметические вычисления)
	power := os.Getenv("COMPUTING_POWER")
	numPower, err := strconv.Atoi(power)

	if err != nil {
		logger.Fatal("Failed to get COMPUTING_POWER value: " + err.Error())
		return
	}

	pool := work.NewPool(numPower, logger)
	defer pool.Shutdown()

	fetcher(pool, logger)
}
