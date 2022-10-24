package service

import (
	"context"
	"fmt"

	pb "github.com/Asliddin3/customer-servis/genproto/customer"
	post "github.com/Asliddin3/customer-servis/genproto/post"

	l "github.com/Asliddin3/customer-servis/pkg/logger"
	grpcclient "github.com/Asliddin3/customer-servis/service/grpc_client"
	"github.com/Asliddin3/customer-servis/storage"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CustomerService struct {
	storage storage.IStorage
	Client  *grpcclient.ServiceManager
	logger  l.Logger
}

func NewCustomerService(clinet *grpcclient.ServiceManager, db *sqlx.DB, log l.Logger) *CustomerService {
	return &CustomerService{
		storage: storage.NewStoragePg(db),
		Client:  clinet,
		logger:  log,
	}
}

func (s *CustomerService) GetCustomerInfo(ctx context.Context, req *pb.CustomerId) (*pb.CustomerResponse, error) {
	customerInfo, err := s.storage.Customer().GetCustomerInfo(req)
	if err != nil {
		s.logger.Error("error getting customer info", l.Any("error getting customer info", err))
		return &pb.CustomerResponse{}, status.Error(codes.Internal, "something went wrong")
	}
	return customerInfo, nil
}
func (s *CustomerService) GetById(ctx context.Context, req *pb.CustomerId) (*pb.CustomerResponsePost, error) {
	customerPosts, err := s.storage.Customer().GetById(req)
	fmt.Print(err)
	if err != nil {
		s.logger.Error("error geting customer", l.Any("error geting custoimer", err))
		return &pb.CustomerResponsePost{}, status.Error(codes.Internal, "something went wrong")
	}
	fmt.Println(customerPosts.Id)
	id := &post.CustomerId{
		Id: customerPosts.Id,
	}
	fmt.Println(id.Id)
	posts, err := s.Client.PostServise().GetPostCustomerId(context.Background(), id)
	// posts,err:=s.storage.Post().GetPostCustomerId(&post.CustomerId{Id: customerPosts.Id})
	fmt.Println(err)
	if err != nil {
		s.logger.Error("error getting post from post-servise", l.Any("error getting post", err))
		return &pb.CustomerResponsePost{}, status.Error(codes.Internal, "something went wrong")
	}
	postsRespose := []*pb.PostResponse{}
	for _, post := range posts.Posts {
		postResp := pb.PostResponse{
			Name:        post.Name,
			Description: post.Description,
			Id:          post.Id,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
		}
		for _, media := range post.Media {
			postResp.Media = append(postResp.Media, &pb.MediasResponse{
				Id:   media.Id,
				Name: media.Name,
				Link: media.Link,
				Type: media.Type,
			})
		}
		postsRespose = append(postsRespose, &postResp)
	}
	customerPosts.Post = postsRespose
	return customerPosts, nil
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, req *pb.CustomerId) (*pb.Empty, error) {
	_, err := s.storage.Customer().DeleteCustomer(req)
	if err != nil {
		s.logger.Error("error while deleteing customer", l.Any("error deleting customer", err))
		return &pb.Empty{}, status.Error(codes.Internal, "something went wrong")
	}
	id := post.CustomerId{Id: req.Id}
	fmt.Println(id)
	_, err = s.Client.PostServise().DeletePostByCustomerId(context.Background(), &id)
	fmt.Println(err)
	if err != nil {
		s.logger.Error("error while deleting post in customerserivse", l.Any("error deleting post", err))
		return &pb.Empty{}, status.Error(codes.Internal, "error deleting post")
	}
	return &pb.Empty{}, nil
}
func (s *CustomerService) GetListCustomers(ctx context.Context, req *pb.Empty) (*pb.ListCustomers, error) {
	customers, err := s.storage.Customer().GetListCustomers(req)
	for _, customer := range customers.ActiveCustomers {
		posts, err := s.Client.PostServise().GetPostCustomerId(ctx, &post.CustomerId{Id: customer.Id})
		if err != nil {
			s.logger.Error("error gettinf customer posts", l.Any("error getting customer post", err))
			return &pb.ListCustomers{}, status.Error(codes.Internal, "something went wrong")
		}
		for _, post := range posts.Posts {
			postResp := pb.PostResponse{
				Id:          post.Id,
				Name:        post.Name,
				Description: post.Description,
				CreatedAt:   post.CreatedAt,
				UpdatedAt:   post.UpdatedAt,
			}
			for _, media := range post.Media {
				mediaResp := pb.MediasResponse{
					Id:   media.Id,
					Link: media.Link,
					Name: media.Name,
					Type: media.Type,
				}
				postResp.Media = append(postResp.Media, &mediaResp)
			}
			customer.Posts = append(customer.Posts, &postResp)
		}
	}
	if err != nil {
		s.logger.Error("error getting customers ", l.Any("error geting customers", err))
		return &pb.ListCustomers{}, status.Error(codes.Internal, "something went wrong")
	}
	return customers, nil
}
func (s *CustomerService) UpdateCustomer(ctx context.Context, req *pb.CustomerResponse) (*pb.CustomerResponse, error) {
	customer, err := s.storage.Customer().UpdateCustomer(req)
	if err != nil {
		s.logger.Error("error while updating customer", l.Any("error updating customer", err))
		return &pb.CustomerResponse{}, status.Error(codes.Internal, "somthing went wrong")
	}
	return customer, nil
}

func (s *CustomerService) CreateCustomer(ctx context.Context, req *pb.CustomerRequest) (*pb.CustomerResponse, error) {
	Customer, err := s.storage.Customer().CreateCustomer(req)
	if err != nil {
		s.logger.Error("error while creating Customer", l.Any("error creating Customer", err))
		return &pb.CustomerResponse{}, status.Error(codes.Internal, "something went wrong")
	}
	return Customer, nil
}
