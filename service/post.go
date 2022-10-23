package service

import (
	"context"

	post "github.com/Asliddin3/customer-servis/genproto/post"
	"github.com/Asliddin3/customer-servis/pkg/logger"
	grpcclient "github.com/Asliddin3/customer-servis/service/grpc_client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostService struct {
	Client *grpcclient.ServiceManager
	Logger logger.Logger
}

func (s *PostService) DeletePostByCustomerId(req *post.CustomerId) (*post.Empty, error) {
	_, err := s.Client.PostServise().DeletePostByCustomerId(context.Background(), req)
	if err != nil {
		s.Logger.Error("error while deleting post", logger.Any("error deleting post", err))
		return &post.Empty{}, status.Error(codes.Internal, "error deleting post")
	}
	return &post.Empty{}, nil
}

func (s *PostService) GetPost(req *post.PostId) (*post.PostResponseCustomer, error) {
	postResp, err := s.Client.PostServise().GetPost(context.Background(), req)
	if err != nil {
		s.Logger.Error("error while getting post", logger.Any("error geting post", err))
		return &post.PostResponseCustomer{}, status.Error(codes.Internal, "please check post id")
	}
	return postResp, nil
}
