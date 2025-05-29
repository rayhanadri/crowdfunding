package repository

import (
	"github.com/rayhanadri/crowdfunding/donation-service/model"
	"github.com/stretchr/testify/mock"
)

type MockUserTransactionInterface interface {
	GetAllTransaction(user_id int) (*[]model.Transaction, error)
	CreateTransaction(user_id int, transaction *model.Transaction) (*model.Transaction, error)
	GetTransactionByID(user_id int, transactionID int) (*model.Transaction, error)
	UpdateTransaction(user_id int, transaction *model.Transaction) (*model.Transaction, error)
	CheckUpdateTransaction(user_id int, transaction *model.Transaction) (*model.Transaction, error)
}

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) GetAllTransaction(user_id int) (*[]model.Transaction, error) {
	args := m.Called(user_id)
	if transactions := args.Get(0); transactions != nil {
		return transactions.(*[]model.Transaction), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTransactionRepository) CreateTransaction(user_id int, transaction *model.Transaction) (*model.Transaction, error) {
	args := m.Called(user_id, transaction)
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
