package test

import (
	"testing"
	"time"

	"github.com/rayhanadri/crowdfunding/donation-service/model"
	"github.com/stretchr/testify/assert"

	"github.com/rayhanadri/crowdfunding/api-gateway/repository"
)

func TestGetTransactionByID_Success(t *testing.T) {
	mockRepo := new(repository.MockTransactionRepository)

	// Representing a mock transaction object
	// mockUserId := 1

	mockTransaction := model.Transaction{
		ID:                 1,
		DonationID:         1,
		InvoiceID:          "INV-12345",
		InvoiceURL:         "https://example.com/invoice/12345",
		InvoiceDescription: "Donation for a cause",
		PaymentMethod:      "EWALLET",
		Amount:             50000,
		Status:             "PENDING",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	mockTransactionPtr := &mockTransaction

	// Representing retrieving a transaction by ID from the database
	mockRepo.On("GetTransactionByID", 1).Return(mockTransactionPtr, nil)
	transactionPtr, err := mockRepo.GetTransactionByID(1)

	// Check if the transaction is retrieved successfully
	assert.NoError(t, err)
	assert.NotNil(t, transactionPtr)

	mockRepo.AssertExpectations(t)
}

func TestGetTransactionByID_Failed(t *testing.T) {
	mockRepo := new(repository.MockTransactionRepository)

	// Representing retrieving a transaction by ID from the database
	mockRepo.On("GetTransactionByID", 1).Return(nil, assert.AnError)
	transactionPtr, err := mockRepo.GetTransactionByID(1)

	// Check if the transaction retrieval failed as expected
	assert.Error(t, err)
	assert.Nil(t, transactionPtr)

	mockRepo.AssertExpectations(t)
}

func TestGetAllTransaction_Success(t *testing.T) {
	mockRepo := new(repository.MockTransactionRepository)

	// Representing transactions retrieved from the database
	// mockUserId := 1
	mockTransactions := []model.Transaction{
		{
			ID:                 1,
			DonationID:         1,
			InvoiceID:          "INV-12345",
			InvoiceURL:         "https://example.com/invoice/12345",
			InvoiceDescription: "Donation for a cause",
			PaymentMethod:      "EWALLET",
			Amount:             50000,
			Status:             "PENDING",
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		},
		{
			ID:                 2,
			DonationID:         2,
			InvoiceID:          "INV-12346",
			InvoiceURL:         "https://example.com/invoice/12346",
			InvoiceDescription: "Donation for a cause",
			PaymentMethod:      "EWALLET",
			Amount:             50000,
			Status:             "PENDING",
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		},
	}

	mockTransactionsPtr := &mockTransactions

	// Representing retrieving all transactions from the database
	mockRepo.On("GetAllTransaction").Return(mockTransactionsPtr, nil)
	transactions, err := mockRepo.GetAllTransaction()

	// Check if the transaction is retrieved successfully
	assert.NoError(t, err)
	assert.NotNil(t, transactions)

	mockRepo.AssertExpectations(t)
}

func TestGetAllTransaction_Failed(t *testing.T) {
	mockRepo := new(repository.MockTransactionRepository)

	// Representing retrieving all transactions from the database
	mockRepo.On("GetAllTransaction").Return(nil, assert.AnError)
	transactions, err := mockRepo.GetAllTransaction()

	// Check if the transaction retrieval failed as expected
	assert.Error(t, err)
	assert.Nil(t, transactions)

	mockRepo.AssertExpectations(t)
}

func TestCreateTransaction_Success(t *testing.T) {
	mockRepo := new(repository.MockTransactionRepository)

	// Representing a transaction created and retrieved from the database
	// mockUserId := 1
	mockTransaction := model.Transaction{
		ID:                 1,
		DonationID:         1,
		InvoiceID:          "INV-12345",
		InvoiceURL:         "https://example.com/invoice/12345",
		InvoiceDescription: "Donation for a cause",
		PaymentMethod:      "EWALLET",
		Amount:             50000,
		Status:             "PENDING",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	mockTransactionPtr := &mockTransaction

	// Representing creating a transaction in the database
	mockRepo.On("CreateTransaction", mockTransactionPtr).Return(mockTransactionPtr, nil)
	transactionPtr, err := mockRepo.CreateTransaction(mockTransactionPtr)

	// Check if the transaction is created successfully
	assert.NoError(t, err)
	assert.NotNil(t, transactionPtr)

	mockRepo.AssertExpectations(t)
}

func TestCreateTransaction_Failed(t *testing.T) {
	mockRepo := new(repository.MockTransactionRepository)

	// Representing creating a transaction in the database
	// mockUserId := 1
	mockTransaction := model.Transaction{
		ID:                 1,
		DonationID:         1,
		InvoiceID:          "INV-12345",
		InvoiceURL:         "https://example.com/invoice/12345",
		InvoiceDescription: "Donation for a cause",
		PaymentMethod:      "EWALLET",
		Amount:             50000,
		Status:             "PENDING",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	mockTransactionPtr := &mockTransaction

	mockRepo.On("CreateTransaction", mockTransactionPtr).Return(nil, assert.AnError)
	transactionPtr, err := mockRepo.CreateTransaction(mockTransactionPtr)

	// Check if the transaction creation failed as expected
	assert.Error(t, err)
	assert.Nil(t, transactionPtr)

	mockRepo.AssertExpectations(t)
}

func TestUpdateTransaction_Success(t *testing.T) {
	mockRepo := new(repository.MockTransactionRepository)

	// Representing a transaction created and retrieved from the database
	// mockUserId := 1
	mockTransaction := model.Transaction{
		ID:                 1,
		DonationID:         1,
		InvoiceID:          "INV-12345",
		InvoiceURL:         "https://example.com/invoice/12345",
		InvoiceDescription: "Donation for a cause",
		PaymentMethod:      "EWALLET",
		Amount:             50000,
		Status:             "PENDING",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	mockTransactionPtr := &mockTransaction

	// Representing updating a transaction in the database
	mockRepo.On("UpdateTransaction", mockTransactionPtr).Return(mockTransactionPtr, nil)
	transactionPtr, err := mockRepo.UpdateTransaction(mockTransactionPtr)

	// Check if the transaction is updated successfully
	assert.NoError(t, err)
	assert.NotNil(t, transactionPtr)

	mockRepo.AssertExpectations(t)
}

func TestUpdateTransaction_Failed(t *testing.T) {
	mockRepo := new(repository.MockTransactionRepository)

	// Representing a transaction created and retrieved from the database
	// mockUserId := 1
	mockTransaction := model.Transaction{
		ID:                 1,
		DonationID:         1,
		InvoiceID:          "INV-12345",
		InvoiceURL:         "https://example.com/invoice/12345",
		InvoiceDescription: "Donation for a cause",
		PaymentMethod:      "EWALLET",
		Amount:             50000,
		Status:             "PENDING",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	mockTransactionPtr := &mockTransaction

	// Representing updating a transaction in the database
	mockRepo.On("UpdateTransaction", mockTransactionPtr).Return(nil, assert.AnError)
	transactionPtr, err := mockRepo.UpdateTransaction(mockTransactionPtr)

	// Check if the transaction update failed as expected
	assert.Error(t, err)
	assert.Nil(t, transactionPtr)

	mockRepo.AssertExpectations(t)
}

func TestSyncTransaction_Success(t *testing.T) {
	mockRepo := new(repository.MockTransactionRepository)

	// Representing a transaction created and retrieved from the database
	// mockUserId := 1
	mockTransaction := model.Transaction{
		ID:                 1,
		DonationID:         1,
		InvoiceID:          "INV-12345",
		InvoiceURL:         "https://example.com/invoice/12345",
		InvoiceDescription: "Donation for a cause",
		PaymentMethod:      "EWALLET",
		Amount:             50000,
		Status:             "PENDING",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	mockTransactionPtr := &mockTransaction

	// Representing syncing a transaction in the database
	mockRepo.On("SyncTransaction", mockTransaction.ID).Return(mockTransactionPtr, nil)
	transactionPtr, err := mockRepo.SyncTransaction(mockTransaction.ID)

	// Check if the transaction is synced successfully
	assert.NoError(t, err)
	assert.NotNil(t, transactionPtr)

	mockRepo.AssertExpectations(t)
}

func TestCheckUpdateTransaction_Failed(t *testing.T) {
	mockRepo := new(repository.MockTransactionRepository)

	// Representing a transaction created and retrieved from the database
	// mockUserId := 1
	mockTransaction := model.Transaction{
		ID:                 1,
		DonationID:         1,
		InvoiceID:          "INV-12345",
		InvoiceURL:         "https://example.com/invoice/12345",
		InvoiceDescription: "Donation for a cause",
		PaymentMethod:      "EWALLET",
		Amount:             50000,
		Status:             "PENDING",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	// mockTransactionPtr := &mockTransaction

	// Representing syncing a transaction in the database
	mockRepo.On("SyncTransaction", mockTransaction.ID).Return(nil, assert.AnError)
	transactionPtr, err := mockRepo.SyncTransaction(mockTransaction.ID)

	// Check if the transaction update failed as expected
	assert.Error(t, err)
	assert.Nil(t, transactionPtr)

	mockRepo.AssertExpectations(t)
}
