package service

import (
	"context"

	"google.golang.org/grpc"

	pb "test_task/proto"
)

type Service struct {
}

func New() (*Service, error) {
	return &Service{}, nil
}

func (s *Service) GetData(ctx context.Context, in *pb.DataRequest, opts ...grpc.CallOption) (*pb.DataResponse, error) {
	// TODO: Write implementation after writing other app logic
	return &pb.DataResponse{Data: nil}, nil
}
