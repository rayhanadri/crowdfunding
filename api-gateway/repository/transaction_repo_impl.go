package repository

import (
	"github.com/rayhanadri/crowdfunding/donation-service/model"
	"github.com/stretchr/testify/mock"
)

type MockUserTransactionInterface interface {
	GetAllTransaction() (*[]model.Transaction, error)
	CreateTransaction(transaction *model.Transaction) (*model.Transaction, error)
	GetTransactionByID(transactionID int) (*model.Transaction, error)
	UpdateTransaction(transaction *model.Transaction) (*model.Transaction, error)
	CheckUpdateTransaction(transaction *model.Transaction) (*model.Transaction, error)
}

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) GetAllTransaction() (*[]model.Transaction, error) {
	args := m.Called()
	if transactions := args.Get(0); transactions != nil {
		return transactions.(*[]model.Transaction), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTransactionRepository) CreateTransaction(transaction *model.Transaction) (*model.Transaction, error) {
	args := m.Called(transaction)
	if transaction := args.Get(0); transaction != nil {
		return transaction.(*model.Transaction), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTransactionRepository) GetTransactionByID(user_id int, transactionID int) (*model.Transaction, error) {
	args := m.Called(user_id, transactionID)
	if transaction := args.Get(0); transaction != nil {
		return transaction.(*model.Transaction), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTransactionRepository) UpdateTransaction(user_id int, transaction *model.Transaction) (*model.Transaction, error) {
	args := m.Called(user_id, transaction)
	if transaction := args.Get(0); transaction != nil {
		return transaction.(*model.Transaction), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTransactionRepository) CheckUpdateTransaction(user_id int, transaction *model.Transaction) (*model.Transaction, error) {
	args := m.Called(user_id, transaction)
	if transaction := args.Get(0); transaction != nil {
		return transaction.(*model.Transaction), args.Error(1)
	}
	return nil, args.Error(1)
}
