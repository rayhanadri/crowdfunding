package service

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"donation-service/config"
	"donation-service/model"
	"donation-service/pb"

)

type DonationService struct {
	pb.UnimplementedDonationServiceServer
}

func (s *DonationService) GetUserByID(ctx context.Context, req *pb.DonationIdRequest) (*pb.DonationResponse, error) {
	// Extract the ID from the request
	id := req.GetId()

	var donation model.Donation
	if err := config.DB.First(&donation, id).Error; err != nil {
		return nil, err
	}

	// Create a donation response
	response := &pb.DonationResponse{
		Id:        int32(donation.ID),
		UserID: int32(donation.UserID),
		CampaignID: int32(donation.CampaignID),
		Amount: float32(donation.amount),
		Message: donation.message,
		Status: donation.status,
		CreatedAt: donation.createdAt,
		UpdatedAt: donation.updatedAt
	}

	return response, nil
}

func (r *DonationService) CreateDonation(ctx context.Context, req *pb.DonationRequest) (*pb.DonationResponse, error) {
	donation := &model.Donation{
		UserID:     int(req.GetId()),
		CampaignID: int(req.GetCampaignId()),
		Amount:     float64(req.GetAmount()),
		Message:    req.GetMessage(),
		Status:     req.GetStatus(),
	}

	//validate user data
	if donation.UserID == 0 || donation.CampaignID == 0 || donation.Amount <= 0 {
		err := errors.New("user ID, campaign ID, and amount are required")
		response := &pb.UserResponse{
			Message: "Failed to create donation",
			Error:   err.Error(),
		}

		return response, err
	}

	if err := config.DB.Omit("id").Create(donation).Error; err != nil {
		response := &pb.UserResponse{
			Message: "Failed to create donation",
			Error:   err.Error(),
		}

		return response, err
	}
	
	if err := config.DB.Last(donation).Error; err != nil {
		response := &pb.UserResponse{
			Message: "Failed to create user",
			Error:   err.Error(),
		}

		return response, err
	}

	// Create a donation response
	response := &pb.DonationResponse{
		Id:        int32(donation.ID),
		UserID:    int32(donation.UserID),
		CampaignID: int32(donation.CampaignID),
		Amount:    float32(donation.Amount),
		Message:   donation.Message,
		Status:    donation.Status,
		CreatedAt: donation.CreatedAt.Format(time.RFC3339),
		UpdatedAt: donation.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (r *DonationService) UpdateDonation(ctx context.Context, req *pb.DonationRequest) (*pb.DonationRequest, error) {
	donation := &model.Donation{
		ID:         int(req.GetId()),
		UserID:     int(req.GetUserId()),
		CampaignID: int(req.GetCampaignId()),
		Amount:     float64(req.GetAmount()),
		Message:    req.GetMessage(),
		Status:     req.GetStatus(),
	}

	if err := config.DB.Model(donation).Updates(donation).Error; err != nil {
		return nil, err
	}

	// Create a user response
	response := &pb.DonationResponse{
		Id:        int32(donation.ID),
		UserID:    int32(donation.UserID),
		CampaignID: int32(donation.CampaignID),
		Amount:    float32(donation.Amount),
		Message:   donation.Message,
		Status:    donation.Status,
		CreatedAt: donation.CreatedAt.Format(time.RFC3339),
		UpdatedAt: donation.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}
