package repository

import (
	"github.com/rayhanadri/crowdfunding/user-service/model"
	"github.com/stretchr/testify/mock"
)

type MockUserRepositoryInterface interface {
	GetUserByID(id int) (*model.User, error)
	CreateUser(user *model.User) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	LoginUser(user *model.User) (*model.User, error)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByID(id int) (*model.User, error) {
	args := m.Called(id)
	if user := args.Get(0); user != nil {
		return user.(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) CreateUser(user *model.User) (*model.User, error) {
	args := m.Called(user)
	if user := args.Get(0); user != nil {
		return user.(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) UpdateUser(user *model.User) (*model.User, error) {
	args := m.Called(user)
	if user := args.Get(0); user != nil {
		return user.(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) LoginUser(user *model.User) (*model.User, error) {
	args := m.Called(user)
	if user := args.Get(0); user != nil {
		return user.(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}
