package repo

import (
	pb "github.com/Asliddin3/customer-servis/genproto/post"
)

type PostStorageI interface {
	GetPost(*pb.PostId) (*pb.PostResponse, error)
	DeletePostByCustomerId(*pb.CustomerId) (*pb.Empty, error)
}
