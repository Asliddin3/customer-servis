package grpcClient

import (
	"fmt"

	"github.com/Asliddin3/customer-servis/config"
	"google.golang.org/grpc"

	postPB "github.com/Asliddin3/customer-servis/genproto/post"
)

//GrpcClientI ...

//GrpcClient ...
type ServiceManager struct {
	conf        config.Config
	postServise postPB.PostServiceClient
}

//New ...
func New(cfg config.Config) (*ServiceManager, error) {
	connPost, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.PostServiceHost, cfg.PostServicePort),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("error while dial post servise: host:%s and port:%d",
			cfg.PostServiceHost, cfg.PostServicePort)
	}
	serviceManager := &ServiceManager{
		conf:        cfg,
		postServise: postPB.NewPostServiceClient(connPost),
	}
	return serviceManager, nil
}

func (s *ServiceManager) PostServise() postPB.PostServiceClient {
	return s.postServise
}
