package repo

import (
	pb "github.com/Asliddin3/customer-servis/genproto/customer"
)

type CustomerStorageI interface {
	CreateCustomer(*pb.CustomerRequest) (*pb.CustomerResponse, error)
	UpdateCustomer(*pb.CustomerUpdate) (*pb.CustomerResponse, error)
	DeleteCustomer(*pb.CustomerId) (*pb.Empty, error)
	GetById(*pb.CustomerId) (*pb.CustomerResponsePost, error)
	GetListCustomers(*pb.Empty) (*pb.ListCustomers, error)
	GetCustomerInfo(*pb.CustomerId) (*pb.CustomerResponse, error)
}
