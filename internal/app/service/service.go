package service

import (
	"context"
	"encoding/json"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
	"net"
	"test_task/internal/app/database"
	pb "test_task/proto"
)

type Service struct {
	pageSize int32
	db       database.IDatabase

	errChan chan error
}

func New(pageSize int32, db database.IDatabase, errChan chan error) (*Service, error) {
	serv := &Service{}

	serv.pageSize = pageSize
	serv.db = db
	serv.errChan = errChan

	return serv, nil
}

func (s *Service) Run() {
	listener, err := net.Listen("tcp", ":5300")
	if err != nil {
		s.errChan <- err
		return
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterApiServiceServer(grpcServer, s)
	err = grpcServer.Serve(listener)
	if err != nil {
		s.errChan <- err
		return
	}
}

func (s *Service) GetData(ctx context.Context, req *pb.DataRequest) (*pb.DataResponse, error) {
	offset := s.pageSize * req.Page
	guid, err := uuid.FromString(req.Guid)
	if err != nil {
		return nil, err
	}

	data, err := s.db.GetDataAPI(ctx, guid, offset, req.Limit)
	if err != nil {
		return nil, err
	}

	var arrSt []*structpb.Struct
	for _, row := range data {
		jsonRow, err := json.Marshal(row)
		if err != nil {
			return nil, err
		}

		st := &structpb.Struct{}
		err = protojson.Unmarshal(jsonRow, st)
		if err != nil {
			return nil, err
		}
		
		arrSt = append(arrSt, st)
	}

	return &pb.DataResponse{Data: arrSt}, nil
}
