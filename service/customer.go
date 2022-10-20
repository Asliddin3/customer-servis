package service

import (
	"context"

	pb "github.com/Asliddin3/customer-servis/genproto/customer"
	l "github.com/Asliddin3/customer-servis/pkg/logger"
	"github.com/Asliddin3/customer-servis/storage"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CustomerService struct {
	storage storage.IStorage
	logger  l.Logger
}

func NewCustomerService(db *sqlx.DB, log l.Logger) *CustomerService {
	return &CustomerService{
		storage: storage.NewStoragePg(db),
		logger:  log,
	}
}



func (s *CustomerService) CreateCustomer(ctx context.Context, req *pb.CustomerRequest) (*pb.CustomerResponse, error) {
	Customer, err := s.storage.Customer().CreateCustomer(req)
	if err != nil {
		s.logger.Error("error while creating Customer", l.Any("error creating Customer", err))
		return &pb.CustomerResponse{}, status.Error(codes.Internal, "something went wrong")
	}
	return Customer, nil
}
