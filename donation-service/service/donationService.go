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

	createdAt, err := time.Parse(time.RFC3339, res.GetCreatedAt())
	if err != nil {
		return nil, fmt.Errorf("error parsing CreatedAt: %v", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, res.GetUpdatedAt())
	if err != nil {
		return nil, fmt.Errorf("error parsing UpdatedAt: %v", err)
	}

	userModel = &user_model.User{
		ID:        int(res.GetId()),
		Name:      res.GetName(),
		Email:     res.GetEmail(),
		Password:  res.GetPassword(),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	return userModel, nil
}

func (s *DonationService) GetAllDonations(ctx context.Context, req *pb.GetDonationsRequest) (*pb.GetDonationsResponse, error) {

	var donations []model.Donation
	if err := config.DB.Find(&donations).Error; err != nil {
		return nil, err
	}

	// Create a donation response
	response := &pb.GetDonationsResponse{
		Donations: make([]*pb.Donation, 0, len(donations)),
	}

	for _, donation := range donations {
		userModel, err := GetUserByID(int32(donation.UserID))
		if err != nil {
			return nil, err
		}
		if userModel == nil {
			return nil, fmt.Errorf("user with ID %d not found", donation.UserID)
		}

		donationResponse := &pb.Donation{
			Id:         int32(donation.ID),
			UserId:     int32(donation.UserID),
			CampaignId: int32(donation.CampaignID),
			Amount:     float32(donation.Amount),
			Message:    donation.MessageText,
			Status:     donation.Status,
			CreatedAt:  donation.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  donation.UpdatedAt.Format(time.RFC3339),
		}
		response.Donations = append(response.Donations, donationResponse)
	}

	return response, nil
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

func (s *DonationService) GetAllTransactions(ctx context.Context, req *pb.GetTransactionRequest) (*pb.GetTransactionResponse, error) {
	var transactions []model.Transaction
	if err := config.DB.Find(&transactions).Error; err != nil {
		return nil, err
	}

	// Create a donation response
	response := &pb.GetTransactionResponse{
		Transactions: make([]*pb.Transaction, 0, len(transactions)),
	}

	// Iterate through each transaction and fetch the associated donation
	for _, transaction := range transactions {
		donation, err := s.GetDonationByID(ctx, &pb.DonationIdRequest{Id: int32(transaction.DonationID)})
		if err != nil {
			return nil, fmt.Errorf("error getting donation for transaction ID %d: %v", transaction.ID, err)
		}

		transactionResponse := &pb.Transaction{
			Id:                 int32(transaction.ID),
			DonationId:         donation.GetId(),
			InvoiceId:          transaction.InvoiceID,
			InvoiceUrl:         transaction.InvoiceURL,
			InvoiceDescription: transaction.InvoiceDescription,
			PaymentMethod:      transaction.PaymentMethod,
			Amount:             float32(transaction.Amount),
			Status:             transaction.Status,
			CreatedAt:          transaction.CreatedAt.Format(time.RFC3339),
			UpdatedAt:          transaction.UpdatedAt.Format(time.RFC3339),
		}
		response.Transactions = append(response.Transactions, transactionResponse)
	}

	return response, nil
}

func (s *DonationService) GetTransactionByID(ctx context.Context, req *pb.TransactionIdRequest) (*pb.TransactionResponse, error) {
	// Extract the ID from the request
	id := req.GetId()

	var transaction model.Transaction
	if err := config.DB.First(&transaction, id).Error; err != nil {
		return nil, err
	}

	// Create a donation response
	response := &pb.TransactionResponse{
		Id:                 int32(transaction.ID),
		DonationId:         int32(transaction.DonationID),
		InvoiceId:          transaction.InvoiceID,
		InvoiceUrl:         transaction.InvoiceURL,
		InvoiceDescription: transaction.InvoiceDescription,
		PaymentMethod:      transaction.PaymentMethod,
		Amount:             float32(transaction.Amount),
		Status:             transaction.Status,
		CreatedAt:          transaction.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          transaction.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (r *DonationService) CreateTransaction(ctx context.Context, req *pb.TransactionRequest) (*pb.TransactionResponse, error) {
	transaction := &model.Transaction{
		ID:                 int(req.GetId()),
		DonationID:         int(req.GetDonationId()),
		InvoiceID:          "",
		InvoiceURL:         "",
		InvoiceDescription: "",
		PaymentMethod:      "",
		Amount:             float64(req.GetAmount()),
		Status:             "",
	}

	//validate user data
	if transaction.DonationID == 0 || transaction.Amount <= 0 {
		err := errors.New("donation ID and amount are required")
		response := &pb.TransactionResponse{
			Message: "Failed to create transaction",
			Error:   err.Error(),
		}

		return response, err
	}

	if err := config.DB.Omit("id").Create(transaction).Error; err != nil {
		response := &pb.TransactionResponse{
			Message: "Failed to create transaction",
			Error:   err.Error(),
		}

		return response, err
	}

	if err := config.DB.Last(transaction).Error; err != nil {
		response := &pb.TransactionResponse{
			Message: "Failed to create transaction",
			Error:   err.Error(),
		}

		return response, err
	}

	// Create a transaction response
	response := &pb.TransactionResponse{
		Message:            "Transaction created successfully",
		Id:                 int32(transaction.ID),
		DonationId:         int32(transaction.DonationID),
		InvoiceId:          transaction.InvoiceID,
		InvoiceUrl:         transaction.InvoiceURL,
		InvoiceDescription: transaction.InvoiceDescription,
		PaymentMethod:      transaction.PaymentMethod,
		Amount:             float32(transaction.Amount),
		Status:             transaction.Status,
		CreatedAt:          transaction.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          transaction.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (r *DonationService) UpdateTransaction(ctx context.Context, req *pb.TransactionRequest) (*pb.TransactionResponse, error) {
	donation := &model.Transaction{
		ID: int(req.GetId()),
	}

	if err := config.DB.Model(donation).Updates(donation).Error; err != nil {
		return nil, err
	}

	// Create a user response
	response := &pb.TransactionResponse{
		Id:                 int32(donation.ID),
		DonationId:         int32(donation.DonationID),
		InvoiceId:          donation.InvoiceID,
		InvoiceUrl:         donation.InvoiceURL,
		InvoiceDescription: donation.InvoiceDescription,
		PaymentMethod:      donation.PaymentMethod,
		Amount:             float32(donation.Amount),
		Status:             donation.Status,
		CreatedAt:          donation.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          donation.UpdatedAt.Format(time.RFC3339),
	}

	return response, nil
}
