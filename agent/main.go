package main

import (
	"agent/config"
	"agent/work"
	"fmt"
	"go.uber.org/zap"
)

func main() {
	// Создание и запуск логгера
	logger, err := zap.NewProduction()

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
	numPower, err := config.GetComputingPower()

	if err != nil {
		logger.Fatal("Failed to get COMPUTING_POWER value: " + err.Error())
	}

	pool := work.NewPool(numPower, logger)
	defer pool.Shutdown()

	fetcher(pool, logger)
}
