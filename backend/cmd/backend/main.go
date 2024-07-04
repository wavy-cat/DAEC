package main

import (
	"fmt"
	"github.com/wavy-cat/DAEC/backend/internal/storage"
	"github.com/wavy-cat/DAEC/backend/internal/tasks"
	"github.com/wavy-cat/DAEC/backend/internal/utils"
	"go.uber.org/zap"
)

func main() {
	// Создание и запуск логгера
	logger, err := zap.NewProduction()

	if err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Println("Error synchronizing logger:", err)
		}
	}(logger)

	logger.Info("Orchestrator is starting")

	// Создание базы данных и очереди
	db := storage.NewStorage[utils.ExpressionData]()
	manager := tasks.NewManager()
	defer manager.ShutdownWatcher()

	// Запуск HTTP сервера
	go startHTTPServer(logger, db, manager)

	// Запуск gRPC сервера
	go startGRPCServer(logger, db, manager)

	select {}
}
