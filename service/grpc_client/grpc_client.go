package grpcClient

import (
	"fmt"

	"github.com/Asliddin3/customer-servis/config"
	"google.golang.org/grpc"

	postPB "github.com/Asliddin3/customer-servis/genproto/post"
	reivewPB "github.com/Asliddin3/customer-servis/genproto/review"
)

//GrpcClientI ...

//GrpcClient ...
type ServiceManager struct {
	conf          config.Config
	postServise   postPB.PostServiceClient
	reviewServise reivewPB.ReviewServiceClient
}

func New(cfg config.Config) (*ServiceManager, error) {
	connPost, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.PostServiceHost, cfg.PostServicePort),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("error while dial post servise: host:%s and port:%d",
			cfg.PostServiceHost, cfg.PostServicePort)
	}
	connReview, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.ReviewServiceHost, cfg.ReviewServicePort),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("error while dial Review servise: host:%s and port:%d",
			cfg.ReviewServiceHost, cfg.ReviewServicePort)
	}
	serviceManager := &ServiceManager{
		conf:          cfg,
		postServise:   postPB.NewPostServiceClient(connPost),
		reviewServise: reivewPB.NewReviewServiceClient(connReview),
	}
	return serviceManager, nil
}

func (s *ServiceManager) ReviewServise() reivewPB.ReviewServiceClient {
	return s.reviewServise
}

func (s *ServiceManager) PostServise() postPB.PostServiceClient {
	return s.postServise
}
