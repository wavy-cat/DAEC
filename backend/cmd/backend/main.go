package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/wavy-cat/DAEC/backend/internal/config"
	"github.com/wavy-cat/DAEC/backend/internal/database"
	"github.com/wavy-cat/DAEC/backend/internal/evaluate"
	"github.com/wavy-cat/DAEC/backend/internal/tasks"
	"go.uber.org/zap"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Создание и запуск логгера
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}
	defer logger.Sync()

	logger.Info("Orchestrator is starting")

	// Создание базы данных
	db, err := sql.Open("sqlite3", config.DatabasePath)
	if err != nil {
		logger.Fatal("error when opening database", zap.String("error", err.Error()))
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		logger.Fatal("impossible to ping the database", zap.String("error", err.Error()))
	}

	if err = database.CreateTables(context.TODO(), db); err != nil {
		logger.Fatal("failed to create required tables in database", zap.String("error", err.Error()))
	}

	// Создание менеджера задач
	manager := tasks.NewManager()
	defer manager.ShutdownWatcher()

	// Восстанавливаем задачи
	err = evaluate.Recovery(db, manager)
	if err != nil {
		logger.Error("failed to restore tasks", zap.String("error", err.Error()))
	}

	// Запуск HTTP сервера
	go startHTTPServer(logger, db, manager)

	// Запуск gRPC сервера
	go startGRPCServer(logger, manager)

	select {}
}
