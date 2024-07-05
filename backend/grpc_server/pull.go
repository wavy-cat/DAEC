package grpc_server

import (
	"context"
	"errors"
	"github.com/wavy-cat/DAEC/backend/internal/utils"
	pb "github.com/wavy-cat/DAEC/backend/proto"
	"go.uber.org/zap"
)

func (s *Server) Pull(_ context.Context, _ *pb.Empty) (*pb.PullTaskResponse, error) {
	task, found := s.manager.GetTask()
	if !found {
		return nil, errors.New("no task yet")
	}

	s.logger.Info("Agent takes on a new task", zap.String("ID", task.Id.String()))
	return &pb.PullTaskResponse{
		Id:            task.Id.String(),
		Arg1:          task.Arg1,
		Arg2:          task.Arg2,
		Operation:     utils.ParseOperation(task.Operation),
		OperationTime: task.OperationTime,
	}, nil
}
