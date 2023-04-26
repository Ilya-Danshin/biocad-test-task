package service

import (
	"context"
	"encoding/json"
	"net"
	"strconv"
	
	"github.com/gofrs/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"

	"test_task/internal/app/config"
	"test_task/internal/app/database"
	pb "test_task/proto"
)

type Service struct {
	pageSize int32
	db       database.IDatabase
	port     int32

	errChan chan error
}

func New(cfg config.Service, db database.IDatabase, errChan chan error) (*Service, error) {
	serv := &Service{}

	serv.pageSize = cfg.PageSize
	serv.db = db
	serv.port = cfg.Port
	serv.errChan = errChan

	return serv, nil
}

func (s *Service) Run() {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(int(s.port)))
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
