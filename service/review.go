package service

import (
	"context"

	review "github.com/Asliddin3/customer-servis/genproto/review"
	"github.com/Asliddin3/customer-servis/pkg/logger"
	grpcclient "github.com/Asliddin3/customer-servis/service/grpc_client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReviewService struct {
	Client *grpcclient.ServiceManager
	Logger logger.Logger
}

func (r *ReviewService) GetCustomerReviews(ctx context.Context, req *review.CustomerId) (*review.CustomerReviewList, error) {
	reviews, err := r.Client.ReviewServise().GetCustomerReviews(context.Background(), req)
	if err != nil {
		r.Logger.Error("error getting customer review in customer", logger.Any("error getting reviews", err))
		return &review.CustomerReviewList{}, status.Error(codes.Internal, "something went wrong")
	}
	return reviews, nil
}
