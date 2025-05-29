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

type TransactionRepository interface {
	GetAllTransaction() (*[]model.Transaction, error)
	CreateTransaction(transaction *model.Transaction) (*model.Transaction, error)
	GetTransactionByID(transactionID int) (*model.Transaction, error)
	UpdateTransaction(transaction *model.Transaction) (*model.Transaction, error)
	SyncTransaction(transactionID int) (*model.Transaction, error)
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
	// Call the GetDonations method
	res, err := client.GetAllTransactions(ctx, req)
	if err != nil {
		log.Printf("Error calling GetAllTransactions: %v", err)
		return nil, err
	}

	var transactions []model.Transaction
	for _, d := range res.GetTransactions() {
		var transaction model.Transaction
		GetCreatedAtTime, err := time.Parse(time.RFC3339, d.GetCreatedAt())
		if err != nil {
			return nil, fmt.Errorf("invalid created_at value: %v", err)
		}
		GetUpdatedAtTime, err := time.Parse(time.RFC3339, d.GetUpdatedAt())
		if err != nil {
			return nil, fmt.Errorf("invalid updated_at value: %v", err)
		}
		transaction.ID = int(d.Id)
		transaction.DonationID = int(d.DonationId)
		transaction.InvoiceID = d.GetInvoiceId()
		transaction.InvoiceURL = d.GetInvoiceUrl()
		transaction.InvoiceDescription = d.GetInvoiceDescription()
		transaction.PaymentMethod = d.GetPaymentMethod()
		transaction.Amount = float64(d.GetAmount())
		transaction.Status = d.GetStatus()
		transaction.CreatedAt = GetCreatedAtTime
		transaction.UpdatedAt = GetUpdatedAtTime

		// get donation from grpc
		donationReq := &pb.DonationIdRequest{Id: int32(transaction.DonationID)}
		donationRes, err := client.GetDonationByID(ctx, donationReq)
		if err != nil {
			log.Printf("Error calling GetDonation: %v", err)
			return nil, err
		}

		GetDonationCreatedAtTime, err := time.Parse(time.RFC3339, donationRes.GetCreatedAt())
		if err != nil {
			return nil, fmt.Errorf("invalid created_at value: %v", err)
		}
		GetDonationUpdatedAtTime, err := time.Parse(time.RFC3339, donationRes.GetUpdatedAt())
		if err != nil {
			return nil, fmt.Errorf("invalid updated_at value: %v", err)
		}

		transaction.Donation = model.Donation{
			ID:          int(donationRes.GetId()),
			CampaignID:  int(donationRes.GetCampaignId()),
			Amount:      float64(donationRes.GetAmount()),
			MessageText: donationRes.GetMessageText(),
			Status:      donationRes.GetStatus(),
			CreatedAt:   GetDonationCreatedAtTime,
			UpdatedAt:   GetDonationUpdatedAtTime,
		}

		// push to arrays
		transactions = append(transactions, transaction)
	}
	if len(transactions) == 0 {
		return nil, errors.New("no transactions found")
	}

	return &transactions, nil
}

func (r *transactionRepository) CreateTransaction(transaction *model.Transaction) (*model.Transaction, error) {
	// call grpc
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
	req := &pb.TransactionRequest{Id: 0, DonationId: int32(transaction.DonationID), InvoiceId: "", InvoiceUrl: "", InvoiceDescription: "", PaymentMethod: "", Amount: float32(transaction.Amount), Status: "PENDING"}
	// Call the CreateTransaction method
	res, err := client.CreateTransaction(ctx, req) // Update to call CreateDonation instead of GetDonationByID
	if err != nil {
		log.Printf("Error calling CreateTransaction: %v", err)
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

	return transaction, nil
}

func (r *transactionRepository) UpdateTransaction(transaction *model.Transaction) (*model.Transaction, error) {
	// call grpc
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
	req := &pb.TransactionRequest{Id: int32(transaction.ID), DonationId: int32(transaction.DonationID), InvoiceId: "", InvoiceUrl: "", InvoiceDescription: "", PaymentMethod: "", Amount: float32(transaction.Amount), Status: "PENDING"}
	// Call the UpdateTransaction method
	res, err := client.UpdateTransaction(ctx, req) // Update to call CreateDonation instead of GetDonationByID
	if err != nil {
		log.Printf("Error calling UpdateTransaction: %v", err)
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

	return transaction, nil
}

func (r *transactionRepository) GetTransactionByID(transactionID int) (*model.Transaction, error) {
	// call grpc
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

func (r *transactionRepository) SyncTransaction(transactionID int) (*model.Transaction, error) {
	// call grpc
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
	res, err := client.SyncTransaction(ctx, req)
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
