package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	user_model "github.com/rayhanadri/crowdfunding/user-service/model"
	user_pb "github.com/rayhanadri/crowdfunding/user-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/rayhanadri/crowdfunding/donation-service/config" // corrected the import path
	"github.com/rayhanadri/crowdfunding/donation-service/model"  // corrected the import path
	"github.com/rayhanadri/crowdfunding/donation-service/pb"     // corrected the import path
)

type DonationService struct {
	pb.UnimplementedDonationServiceServer
}

func GetUserByID(userId int32) (userModel *user_model.User, error error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new client
	client := user_pb.NewUserServiceClient(conn)
	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Create a request
	req := &user_pb.UserIdRequest{Id: userId}
	// Call the GetUserByID method
	res, err := client.GetUserByID(ctx, req)
	if err != nil {
		log.Fatalf("Error calling GetUserByID: %v", err)
		return nil, err
	}

	userModel = user_model.User{
		ID:        int(res.GetId()),
		Name:      res.GetName(),
		Email:     res.GetEmail(),
		Password:  res.GetPassword(),
		CreatedAt: res.GetCreatedAt(),
		UpdatedAt: res.GetUpdatedAt(),
	}
	return &userModel, nil
}

func (s *DonationService) GetDonationByID(ctx context.Context, req *pb.DonationIdRequest) (*pb.DonationResponse, error) {
	// Extract the ID from the request
	id := req.GetId()

	var donation model.Donation
	if err := config.DB.First(&donation, id).Error; err != nil {
		return nil, err
	}

	userModel, err := GetUserByID(int32(donation.UserID))
	if err != nil {
		return nil, err
	}
	if userModel == nil {
		return nil, fmt.Errorf("user with ID %d not found", donation.UserID)
	}

	// Create a donation response
	response := &pb.DonationResponse{
		Id:          int32(donation.ID),
		UserId:      int32(donation.UserID),
		CampaignId:  int32(donation.CampaignID),
		Amount:      float32(donation.Amount),
		MessageText: donation.MessageText,
		Status:      donation.Status,
		CreatedAt:   donation.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   donation.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (r *DonationService) CreateDonation(ctx context.Context, req *pb.DonationRequest) (*pb.DonationResponse, error) {
	donation := &model.Donation{
		UserID:      int(req.GetId()),
		CampaignID:  int(req.GetCampaignId()),
		Amount:      float64(req.GetAmount()),
		MessageText: req.GetMessage(),
		Status:      req.GetStatus(),
	}

	//validate user data
	if donation.UserID == 0 || donation.CampaignID == 0 || donation.Amount <= 0 {
		err := errors.New("user ID, campaign ID, and amount are required")
		response := &pb.DonationResponse{
			Message: "Failed to create donation",
			Error:   err.Error(),
		}

		return response, err
	}

	if err := config.DB.Omit("id").Create(donation).Error; err != nil {
		response := &pb.DonationResponse{
			Message: "Failed to create donation",
			Error:   err.Error(),
		}

		return response, err
	}

	if err := config.DB.Last(donation).Error; err != nil {
		response := &pb.DonationResponse{
			Message: "Failed to create user",
			Error:   err.Error(),
		}

		return response, err
	}

	// Create a donation response
	response := &pb.DonationResponse{
		Message:     "Donation created successfully",
		Id:          int32(donation.ID),
		UserId:      int32(donation.UserID),
		CampaignId:  int32(donation.CampaignID),
		Amount:      float32(donation.Amount),
		MessageText: donation.MessageText,
		Status:      donation.Status,
		CreatedAt:   donation.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   donation.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (r *DonationService) UpdateDonation(ctx context.Context, req *pb.DonationRequest) (*pb.DonationResponse, error) {
	donation := &model.Donation{
		ID:          int(req.GetId()),
		UserID:      int(req.GetUserId()),
		CampaignID:  int(req.GetCampaignId()),
		Amount:      float64(req.GetAmount()),
		MessageText: req.GetMessage(),
		Status:      req.GetStatus(),
	}

	if err := config.DB.Model(donation).Updates(donation).Error; err != nil {
		return nil, err
	}

	// Create a user response
	response := &pb.DonationResponse{
		Id:         int32(donation.ID),
		UserId:     int32(donation.UserID),
		CampaignId: int32(donation.CampaignID),
		Amount:     float32(donation.Amount),
		Message:    donation.MessageText,
		Status:     donation.Status,
		CreatedAt:  donation.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  donation.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}
