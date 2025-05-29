package test

import (
	"testing"
	"time"

	"github.com/rayhanadri/crowdfunding/user-service/model"
	"github.com/stretchr/testify/assert"

	"github.com/rayhanadri/crowdfunding/api-gateway/repository"
)

func TestGetUserByID_Success(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)

	// Representing a user retrieved from the database
	mockUser := model.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockUserPtr := &mockUser

	// Representing retrieving a user by ID from the database
	mockRepo.On("GetUserByID", 1).Return(mockUserPtr, nil)
	userPtr, err := mockRepo.GetUserByID(1)

	// Check if the user is retrieved successfully
	assert.NoError(t, err)
	assert.NotNil(t, userPtr)

	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_Failed(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)

	// Representing retrieving a user by ID from the database
	mockRepo.On("GetUserByID", 2).Return(nil, assert.AnError)
	userPtr, err := mockRepo.GetUserByID(2)

	// Check if the user retrieval failed as expected
	assert.Error(t, err)
	assert.Nil(t, userPtr)

	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)

	// Representing a user created and retrieved from the database
	mockUser := model.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockUserPtr := &mockUser

	// Representing creating a user in the database
	mockRepo.On("CreateUser", mockUserPtr).Return(mockUserPtr, nil)
	userPtr, err := mockRepo.CreateUser(mockUserPtr)

	// Check if the user is created successfully
	assert.NoError(t, err)
	assert.NotNil(t, userPtr)

	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Failed(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)

	// Representing a user to be created
	mockUser := model.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockUserPtr := &mockUser

	// Representing creating a user in the database
	mockRepo.On("CreateUser", mockUserPtr).Return(nil, assert.AnError)
	userPtr, err := mockRepo.CreateUser(mockUserPtr)

	// Check if the user creation failed as expected
	assert.Error(t, err)
	assert.Nil(t, userPtr)

	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)

	// Representing a user to be updated
	mockUser := model.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockUserPtr := &mockUser

	// Representing updating a user in the database
	mockRepo.On("UpdateUser", mockUserPtr).Return(mockUserPtr, nil)
	userPtr, err := mockRepo.UpdateUser(mockUserPtr)

	// Check if the user is updated successfully
	assert.NoError(t, err)
	assert.NotNil(t, userPtr)

	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_Failed(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)

	// Representing a user to be updated
	mockUser := model.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockUserPtr := &mockUser

	// Representing updating a user in the database
	mockRepo.On("UpdateUser", mockUserPtr).Return(nil, assert.AnError)
	userPtr, err := mockRepo.UpdateUser(mockUserPtr)

	// Check if the user update failed as expected
	assert.Error(t, err)
	assert.Nil(t, userPtr)

	mockRepo.AssertExpectations(t)
}

func TestLoginUser_Success(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)

	// Representing a user to be logged in
	mockUser := model.User{
		Email:    "john.doe@example.com",
		Password: "password123",
	}
	mockUserPtr := &mockUser

	// Representing logging in a user in the database
	mockRepo.On("LoginUser", mockUserPtr).Return(mockUserPtr, nil)
	userPtr, err := mockRepo.LoginUser(mockUserPtr)

	// Check if the user login is successful
	assert.NoError(t, err)
	assert.NotNil(t, userPtr)

	mockRepo.AssertExpectations(t)
}

func TestLoginUser_Failed(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)

	// Representing a user to be logged in
	mockUser := model.User{
		Email:    "john.doe@example.com",
		Password: "password123",
	}
	mockUserPtr := &mockUser

	// Representing logging in a user in the database
	mockRepo.On("LoginUser", mockUserPtr).Return(nil, assert.AnError)
	userPtr, err := mockRepo.LoginUser(mockUserPtr)

	// Check if the user login failed as expected
	assert.Error(t, err)
	assert.Nil(t, userPtr)

	mockRepo.AssertExpectations(t)
}
