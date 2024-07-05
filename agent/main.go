package main

import (
	"fmt"
	"github.com/wavy-cat/DAEC/agent/config"
	pb "github.com/wavy-cat/DAEC/agent/proto"
	"github.com/wavy-cat/DAEC/agent/work"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	// Соединение с оркестратором
	conn, err := grpc.NewClient(config.BackendAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}(conn)
	grpcClient := pb.NewTasksServiceClient(conn)

	// Запуск worker pool (горутин, которые будут выполнять арифметические вычисления)
	numPower, err := config.GetComputingPower()

	if err != nil {
		logger.Fatal("Failed to get COMPUTING_POWER value: " + err.Error())
	}

	pool := work.NewPool(numPower, logger, grpcClient)
	defer pool.Shutdown()

	fetcher(pool, logger, grpcClient)
}
