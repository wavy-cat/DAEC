package grpc_server

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/wavy-cat/DAEC/backend/proto"
	"go.uber.org/zap"
)

func (s *Server) Push(_ context.Context, in *pb.PushTaskRequest) (*pb.Empty, error) {
	id, err := uuid.Parse(in.Id)
	if err != nil {
		s.logger.Warn(err.Error())
		return nil, err
	}

	err = s.manager.AddResultToTask(id, in.Result, in.Successful)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	s.logger.Info("The agent returned the result of the task", zap.String("ID", in.Id))
	return &pb.Empty{}, nil
}
