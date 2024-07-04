package main

import (
	"github.com/wavy-cat/DAEC/backend/internal/storage"
	"github.com/wavy-cat/DAEC/backend/internal/tasks"
	"github.com/wavy-cat/DAEC/backend/internal/utils"
	"go.uber.org/zap"
)

func startGRPCServer(logger *zap.Logger, db *storage.Storage[utils.ExpressionData], manager *tasks.Manager) {
	logger.Info("Starting gRPC server...")
}
