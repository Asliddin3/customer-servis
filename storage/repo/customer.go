package repo

import (
	pb "github.com/Asliddin3/customer-servis/genproto/customer"
)

type CustomerStorageI interface {
	// CheckField(*pb.CheckFieldRequest) (*pb.CheckFieldResponse,error)
	CreateCustomer(*pb.CustomerRequest)(*pb.CustomerResponse,error)
}
