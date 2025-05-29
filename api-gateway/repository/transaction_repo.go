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

	"github.com/rayhanadri/crowdfunding/api-gateway/entity"
)

type TransactionRepository interface {
	GetAllTransaction() (*[]model.Transaction, error)
	CreateTransaction(transaction *model.Transaction) (*model.Transaction, error)
	GetTransactionByID(transactionID int) (*model.Transaction, error)
	UpdateTransaction(transaction *model.Transaction) (*model.Transaction, error)
	CheckUpdateTransaction(transaction *model.Transaction) (*model.Transaction, error)
}

type transactionRepository struct {
	address string
}

func NewTransactionRepository(address string) TransactionRepository {
	return &transactionRepository{address: address}
}

func (r *transactionRepository) GetAllTransaction() (*[]model.Transaction, error) {
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
	req := &pb.GetTransactionsRequest{} // Use the provided donationID parameter
	// Call the GetTransactions method
	res, err := client.GetAllTransactionss(ctx, req)
	if err != nil {
		log.Printf("Error calling GetAllTransactions: %v", err)
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

func (r *transactionRepository) CreateTransaction(user_id int, transaction *model.Transaction) (*model.Transaction, error) {
	// validate user id
	user := new(entity.User)
	if err := r.db.Where("id = ?", user_id).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// get donation from donation id
	donation := new(entity.Donation)
	if err := r.db.Where("id = ?", transaction.DonationID).First(&donation).Error; err != nil {
		return nil, errors.New("donation not found")
	}

	// check if donation user_id match
	if donation.UserID != user_id {
		return nil, errors.New("user not authorized")
	}

	// check if donation with user and donation id already exist
	var existingTransaction model.Transaction
	if err := r.db.Where("user_id = ? AND donation_id = ?", user_id, transaction.DonationID).First(&existingTransaction).Error; err == nil {
		return nil, errors.New("transaction already exists")
	}

	// get campaign from donation
	campaign := new(entity.Campaign)
	if err := r.db.Where("id = ?", donation.CampaignID).First(&campaign).Error; err != nil {
		return nil, errors.New("campaign not found")
	}
	// get campaign creator user
	userCreator := new(entity.User)
	if err := r.db.Where("id = ?", campaign.UserID).First(&userCreator).Error; err != nil {
		return nil, errors.New("user not found")
	}

	//fill transaction description based on campaign title
	transaction.InvoiceDescription = campaign.Title

	//query create db
	if err := r.db.Create(transaction).Error; err != nil {
		return nil, err
	}

	// Assign All Inner Object for json
	campaign.User = *userCreator
	donation.Campaign = *campaign
	donation.User = *user
	transaction.Donation = model.Donation{
		ID:         donation.ID,
		CampaignID: donation.CampaignID,
		UserID:     donation.UserID,
		Amount:     donation.Amount,
		Status:     donation.Status,
		CreatedAt:  donation.CreatedAt,
		UpdatedAt:  donation.UpdatedAt,
	}

	return transaction, nil
}

func (r *transactionRepository) UpdateTransaction(user_id int, transaction *model.Transaction) (*model.Transaction, error) {
	// validate user id
	user := new(entity.User)
	if err := r.db.Where("id = ?", user_id).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// get donation from donation id
	donation := new(entity.Donation)
	if err := r.db.Where("id = ?", transaction.DonationID).First(&donation).Error; err != nil {
		return nil, errors.New("donation not found")
	}

	// check if donation user_id match
	if donation.UserID != user_id {
		return nil, errors.New("user not authorized")
	}

	// get campaign from donation
	campaign := new(entity.Campaign)
	if err := r.db.Where("id = ?", donation.CampaignID).First(&campaign).Error; err != nil {
		return nil, errors.New("campaign not found")
	}
	// get campaign creator user
	userCreator := new(entity.User)
	if err := r.db.Where("id = ?", campaign.UserID).First(&userCreator).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// get current date as update date
	transaction.UpdatedAt = time.Now()

	// query save transaction
	if err := r.db.Save(transaction).Error; err != nil {
		return nil, err
	}

	// Assign All Inner Object for json
	campaign.User = *userCreator
	donation.Campaign = *campaign
	transaction.Donation = model.Donation{
		ID:         donation.ID,
		CampaignID: donation.CampaignID,
		UserID:     donation.UserID,
		Amount:     donation.Amount,
		Status:     donation.Status,
		CreatedAt:  donation.CreatedAt,
		UpdatedAt:  donation.UpdatedAt,
	}

	return transaction, nil
}

func (r *transactionRepository) GetTransactionByID(transactionID int) (*model.Transaction, error) {
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
	req := &pb.TransactionIdRequest{Id: int32(transactionID)} // Use the provided transactionID parameter
	// Call the GetTransactionByID method
	res, err := client.GetTransactionByID(ctx, req)
	if err != nil {
		log.Printf("Error calling GetTransactionByID: %v", err)
		return nil, err
	}

	var transaction model.Transaction
	GetCreatedAtTime, err := time.Parse(time.RFC3339, res.GetCreatedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid created_at value: %v", err)
	}
	GetUpdatedAtTime, err := time.Parse(time.RFC3339, res.GetUpdatedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid updated_at value: %v", err)
	}

	transaction.ID = int(res.Id)
	transaction.DonationID = int(res.DonationId)
	transaction.InvoiceID = res.GetInvoiceId()
	transaction.InvoiceURL = res.GetInvoiceUrl()
	transaction.InvoiceDescription = res.GetInvoiceDescription()
	transaction.PaymentMethod = res.GetPaymentMethod()
	transaction.Amount = float64(res.GetAmount())
	transaction.Status = res.GetStatus()
	transaction.CreatedAt = GetCreatedAtTime
	transaction.UpdatedAt = GetUpdatedAtTime

	return &transaction, nil
}

func (r *transactionRepository) CheckUpdateTransaction(user_id int, transaction *model.Transaction) (*model.Transaction, error) {
	// validate user id
	user := new(entity.User)
	if err := r.db.Where("id = ?", user_id).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// get donation from donation id
	donation := new(entity.Donation)
	if err := r.db.Where("id = ?", transaction.DonationID).First(&donation).Error; err != nil {
		return nil, errors.New("donation not found")
	}

	// check if donation user_id match
	if donation.UserID != user_id {
		return nil, errors.New("user not authorized")
	}

	// get campaign from donation
	campaign := new(entity.Campaign)
	if err := r.db.Where("id = ?", donation.CampaignID).First(&campaign).Error; err != nil {
		return nil, errors.New("campaign not found")
	}
	// get campaign creator user
	userCreator := new(entity.User)
	if err := r.db.Where("id = ?", campaign.UserID).First(&userCreator).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// get current date as update date
	transaction.UpdatedAt = time.Now()

	// query save transaction
	if err := r.db.Save(transaction).Error; err != nil {
		return nil, err
	}

	// update donation
	donation.Status = transaction.Status

	// update campaign
	if transaction.Status == "PAID" {
		campaign.CollectedAmount += transaction.Amount

		// save campaign
		if err := r.db.Save(campaign).Error; err != nil {
			return nil, err
		}
	}

	// save donation
	if err := r.db.Save(donation).Error; err != nil {
		return nil, err
	}

	// Assign All Inner Object for json
	campaign.User = *userCreator
	donation.Campaign = *campaign
	transaction.Donation = model.Donation{
		ID:         donation.ID,
		CampaignID: donation.CampaignID,
		UserID:     donation.UserID,
		Amount:     donation.Amount,
		Status:     donation.Status,
		CreatedAt:  donation.CreatedAt,
		UpdatedAt:  donation.UpdatedAt,
	}

	return transaction, nil
}
