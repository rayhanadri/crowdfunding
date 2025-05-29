package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/rayhanadri/crowdfunding/donation-service/model"
	"github.com/rayhanadri/crowdfunding/donation-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type DonationRepository interface {
	GetAllDonations() (*[]model.Donation, error)
	CreateDonation(donation *model.Donation) (*model.Donation, error)
	GetDonationByID(donationID int) (*model.Donation, error)
	UpdateDonation(donation *model.Donation) (*model.Donation, error)
}

type donationRepository struct {
	address string
}

func NewDonationRepository(address string) DonationRepository {
	return &donationRepository{address: address}
}

func (r *donationRepository) GetAllDonations() (*[]model.Donation, error) {
	// validate user id
	conn, err := grpc.Dial(
		r.address,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")), // for secure TLS
	)

	if err != nil {
		log.Printf("Did not connect: %v", err)
		return nil, err
	}

	defer conn.Close()

	// Create a new client
	client := pb.NewDonationServiceClient(conn)
	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a request
	req := &pb.GetDonationsRequest{} // Use the provided donationID parameter
	// Call the GetDonations method
	res, err := client.GetAllDonations(ctx, req)
	if err != nil {
		log.Printf("Error calling GetAllDonations: %v", err)
		return nil, err
	}

	var donations []model.Donation
	for _, d := range res.GetDonations() {
		var donation model.Donation
		GetCreatedAtTime, err := time.Parse(time.RFC3339, d.GetCreatedAt())
		if err != nil {
			return nil, fmt.Errorf("invalid created_at value: %v", err)
		}
		GetUpdatedAtTime, err := time.Parse(time.RFC3339, d.GetUpdatedAt())
		if err != nil {
			return nil, fmt.Errorf("invalid updated_at value: %v", err)
		}
		donation.ID = int(d.Id)
		donation.UserID = int(d.UserId)
		donation.CampaignID = int(d.CampaignId)
		donation.Amount = float64(d.GetAmount())
		donation.MessageText = d.GetMessage()
		donation.Status = d.GetStatus()
		donation.CreatedAt = GetCreatedAtTime
		donation.UpdatedAt = GetUpdatedAtTime
		donations = append(donations, donation)
	}
	if len(donations) == 0 {
		return nil, errors.New("no donations found")
	}

	return &donations, nil
}

func (r *donationRepository) GetDonationByID(donationID int) (*model.Donation, error) {
	conn, err := grpc.Dial(
		r.address,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")), // for secure TLS
	)

	if err != nil {
		log.Printf("Did not connect: %v", err)
		return nil, err
	}

	defer conn.Close()

	// Create a new client
	client := pb.NewDonationServiceClient(conn)
	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a request
	req := &pb.DonationIdRequest{Id: int32(donationID)} // Use the provided donationID parameter
	// Call the GetDonationByID method
	res, err := client.GetDonationByID(ctx, req)
	if err != nil {
		log.Printf("Error calling GetDonationByID: %v", err)
		return nil, err
	}

	var donation model.Donation
	GetCreatedAtTime, err := time.Parse(time.RFC3339, res.GetCreatedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid created_at value: %v", err)
	}
	GetUpdatedAtTime, err := time.Parse(time.RFC3339, res.GetUpdatedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid updated_at value: %v", err)
	}

	donation.ID = int(res.Id)
	donation.UserID = int(res.UserId)
	donation.CampaignID = int(res.CampaignId)
	donation.Amount = float64(res.GetAmount())
	donation.MessageText = res.GetMessage()
	donation.Status = res.GetStatus()
	donation.CreatedAt = GetCreatedAtTime
	donation.UpdatedAt = GetUpdatedAtTime

	if donationID == 0 {
		return nil, fmt.Errorf("donation with id %d not found", donationID)
	}

	return &donation, nil
}

func (r *donationRepository) CreateDonation(donation *model.Donation) (*model.Donation, error) {
	// validate user id
	conn, err := grpc.Dial(
		r.address,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")), // for secure TLS
	)

	if err != nil {
		log.Printf("Did not connect: %v", err)
		return nil, err
	}

	defer conn.Close()

	// Create a new client
	client := pb.NewDonationServiceClient(conn)
	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a request
	req := &pb.DonationRequest{Id: int32(donation.ID), UserId: int32(donation.UserID), CampaignId: int32(donation.CampaignID), Amount: float32(donation.Amount), Message: donation.MessageText, Status: donation.Status} // Use the provided donation parameter
	// Call the CreateDonation method
	res, err := client.CreateDonation(ctx, req) // Update to call CreateDonation instead of GetDonationByID
	if err != nil {
		log.Printf("Error calling CreateDonation: %v", err)
		return nil, err
	}

	GetCreatedAtTime, err := time.Parse(time.RFC3339, res.GetCreatedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid created_at value: %v", err)
	}
	GetUpdatedAtTime, err := time.Parse(time.RFC3339, res.GetUpdatedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid updated_at value: %v", err)
	}

	donation.ID = int(res.Id)
	donation.UserID = int(res.UserId)
	donation.CampaignID = int(res.CampaignId)
	donation.Amount = float64(res.GetAmount())
	donation.MessageText = res.GetMessageText()
	donation.Status = res.GetStatus()
	donation.CreatedAt = GetCreatedAtTime
	donation.UpdatedAt = GetUpdatedAtTime

	return donation, nil
}

func (r *donationRepository) UpdateDonation(donation *model.Donation) (*model.Donation, error) {
	// validate user id
	conn, err := grpc.Dial(
		r.address,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")), // for secure TLS
	)

	if err != nil {
		log.Printf("Did not connect: %v", err)
		return nil, err
	}

	defer conn.Close()

	// Create a new client
	client := pb.NewDonationServiceClient(conn)
	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a request
	req := &pb.DonationRequest{Id: int32(donation.ID), UserId: int32(donation.UserID), CampaignId: int32(donation.CampaignID), Amount: float32(donation.Amount), Message: donation.MessageText, Status: donation.Status} // Use the provided donationID parameter
	// Call the CreateDonation method
	res, err := client.UpdateDonation(ctx, req) // Update to call CreateDonation instead of GetDonationByID
	if err != nil {
		log.Printf("Error calling UpdateDonation: %v", err)
		return nil, err
	}

	GetCreatedAtTime, err := time.Parse(time.RFC3339, res.GetCreatedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid created_at value: %v", err)
	}
	GetUpdatedAtTime, err := time.Parse(time.RFC3339, res.GetUpdatedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid updated_at value: %v", err)
	}

	donation.ID = int(res.Id)
	donation.UserID = int(res.UserId)
	donation.CampaignID = int(res.CampaignId)
	donation.Amount = float64(res.GetAmount())
	donation.MessageText = res.GetMessageText()
	donation.Status = res.GetStatus()
	donation.CreatedAt = GetCreatedAtTime
	donation.UpdatedAt = GetUpdatedAtTime

	return donation, nil
}
