package grpc_server

import (
	"github.com/wavy-cat/DAEC/backend/internal/tasks"
	pb "github.com/wavy-cat/DAEC/backend/proto"
	"go.uber.org/zap"
)

type Server struct {
	pb.TasksServiceServer
	manager *tasks.Manager
	logger  *zap.Logger
}

func NewServer(manager *tasks.Manager, logger *zap.Logger) *Server {
	return &Server{
		manager: manager,
		logger:  logger,
	}
}
