package main

import (
	"github.com/wavy-cat/DAEC/backend/grpc_server"
	"github.com/wavy-cat/DAEC/backend/internal/config"
	"github.com/wavy-cat/DAEC/backend/internal/tasks"
	pb "github.com/wavy-cat/DAEC/backend/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func startGRPCServer(logger *zap.Logger, manager *tasks.Manager) {
	logger.Info("Starting gRPC server...")

	lis, err := net.Listen("tcp", config.GRPCAddress)
	if err != nil {
		logger.Fatal(err.Error())
	}

	grpcServer := grpc.NewServer()
	serviceServer := grpc_server.NewServer(manager, logger)
	pb.RegisterTasksServiceServer(grpcServer, serviceServer)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal(err.Error())
	}
}
