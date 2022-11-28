package service

import (
	"context"
	"fmt"

	pb "github.com/Asliddin3/customer-servis/genproto/customer"
	pbp "github.com/Asliddin3/customer-servis/genproto/post"
	post "github.com/Asliddin3/customer-servis/genproto/post"
	"github.com/Asliddin3/customer-servis/genproto/review"

	l "github.com/Asliddin3/customer-servis/pkg/logger"
	"github.com/Asliddin3/customer-servis/pkg/messagebroker"
	grpcclient "github.com/Asliddin3/customer-servis/service/grpc_client"
	"github.com/Asliddin3/customer-servis/storage"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CustomerService struct {
	storage  storage.IStorage
	Client   *grpcclient.ServiceManager
	logger   l.Logger
	producer map[string]messagebroker.Producer
}

func NewCustomerService(clinet *grpcclient.ServiceManager, db *sqlx.DB, log l.Logger, publisher map[string]messagebroker.Producer) *CustomerService {
	return &CustomerService{
		storage:  storage.NewStoragePg(db),
		Client:   clinet,
		logger:   log,
		producer: publisher,
	}
}

func (s *CustomerService) produceMessage(raw *pb.CustomerResponse) error {
	// data, err := raw.Marshal()
	// if err != nil {
	// 	return err
	// }
	post := pbp.PostRequest{
		Name:        "Asliddin",
		Description: "something ",
		CustomerId:  raw.Id,
	}
	data, err := post.Marshal()
	if err != nil {
		return err
	}
	logPost := post.String()
	fmt.Println(data, logPost)
	err = s.producer["customer"].Produce([]byte("customer"), data, logPost)
	fmt.Println("producer error in func", err)
	if err != nil {
		return err
	}
	return nil
}

func (s *CustomerService) CreateCustomer(ctx context.Context, req *pb.CustomerRequest) (*pb.CustomerResponse, error) {
	Customer, err := s.storage.Customer().CreateCustomer(req)
	fmt.Println("servis", Customer, err)
	if err != nil {
		s.logger.Error("error while creating Customer", l.Any("error creating Customer", err))
		return &pb.CustomerResponse{}, status.Error(codes.Internal, "something went wrong")
	}

	err = s.produceMessage(Customer)
	fmt.Println("produce err", err)
	fmt.Println(`
	----------------------------------
	----------------------------------\n`, err)
	if err != nil {
		s.logger.Error("Error while produce to Kafka", l.Any("error produce customer", err))
		return &pb.CustomerResponse{}, status.Error(codes.Internal, err.Error())
	}
	return Customer, nil
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
		reviews, err := s.Client.ReviewServise().GetCustomerReviews(context.Background(), &review.CustomerId{Id: customer.Id})
		if err != nil {
			s.logger.Error("error getting customer reviews ", l.Any("error getting reviews", err))
			return &pb.ListCustomers{}, status.Error(codes.Internal, "something went wrong")
		}
		for _, review := range reviews.ReviewList {
			customer.Reviews = append(customer.Reviews, &pb.ReviewList{
				Id:          review.Id,
				PostId:      review.PostId,
				Description: review.Description,
				Review:      review.Review,
			})
		}
	}
	if err != nil {
		s.logger.Error("error getting customers ", l.Any("error geting customers", err))
		return &pb.ListCustomers{}, status.Error(codes.Internal, "something went wrong")
	}
	return customers, nil
}
func (s *CustomerService) UpdateCustomer(ctx context.Context, req *pb.CustomerUpdate) (*pb.CustomerResponse, error) {
	customer, err := s.storage.Customer().UpdateCustomer(req)
	if err != nil {
		s.logger.Error("error while updating customer", l.Any("error updating customer", err))
		return &pb.CustomerResponse{}, status.Error(codes.Internal, "somthing went wrong")
	}
	return customer, nil
}

func (s *CustomerService) CheckField(ctx context.Context, req *pb.CheckFieldRequest) (*pb.CheckFieldResponse, error) {
	exist, err := s.storage.Customer().CheckField(req)
	if err != nil {
		s.logger.Error("error checking customer", l.Any("error checking customer", err))
		return &pb.CheckFieldResponse{}, status.Error(codes.Internal, "something went wrong")
	}
	return exist, nil
}

func (s *CustomerService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	loginResp, err := s.storage.Customer().Login(req)
	if err != nil {
		s.logger.Error("error while logining", l.Error(err))
		return &pb.LoginResponse{}, status.Error(codes.InvalidArgument, "wrong username or password")
	}
	return loginResp, nil
}

func (s *CustomerService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.LoginResponse, error) {
	refreshToken, err := s.storage.Customer().RefreshToken(req)
	if err != nil {
		s.logger.Error("error refreshing token", l.Any("error updating access token", err))
		return &pb.LoginResponse{}, err
	}
	return refreshToken, nil
}
