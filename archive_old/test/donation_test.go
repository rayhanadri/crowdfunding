package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"crowdfund/model"
	"crowdfund/repository"
)

func TestGetDonationByID_Success(t *testing.T) {
	mockRepo := new(repository.MockDonationRepository)

	// Representing a mock donation object
	mockUserId := 1

	mockDonation := model.Donation{
		ID:         1,
		UserID:     1,
		CampaignID: 1,
		Amount:     50000,
		Message:    "Donation for a cause",
		Status:     "PENDING",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	mockDonationPtr := &mockDonation

	// Representing retrieving a donation by ID from the database
	mockRepo.On("GetDonationByID", mockUserId, 1).Return(mockDonationPtr, nil)
	donationPtr, err := mockRepo.GetDonationByID(mockUserId, 1)

	// Check if the donation is retrieved successfully
	assert.NoError(t, err)
	assert.NotNil(t, donationPtr)

	mockRepo.AssertExpectations(t)
}

func TestGetDonationByID_Failed(t *testing.T) {
	mockRepo := new(repository.MockDonationRepository)

	// Representing retrieving a donation by ID from the database
	mockUserId := 1
	mockRepo.On("GetDonationByID", mockUserId, 1).Return(nil, assert.AnError)
	donationPtr, err := mockRepo.GetDonationByID(mockUserId, 1)

	// Check if the donation retrieval failed as expected
	assert.Error(t, err)
	assert.Nil(t, donationPtr)

	mockRepo.AssertExpectations(t)
}

func TestGetAllDonations_Success(t *testing.T) {
	mockRepo := new(repository.MockDonationRepository)

	// Representing donations retrieved from the database
	mockUserId := 1
	mockDonations := []model.Donation{
		{
			ID:         1,
			UserID:     1,
			CampaignID: 1,
			Amount:     50000,
			Message:    "Donation for a cause",
			Status:     "PENDING",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         2,
			UserID:     2,
			CampaignID: 2,
			Amount:     50000,
			Message:    "Donation for a cause",
			Status:     "PENDING",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	mockDonationsPtr := &mockDonations

	// Representing retrieving all donations from the database
	mockRepo.On("GetAllDonation", mockUserId).Return(mockDonationsPtr, nil)
	donations, err := mockRepo.GetAllDonation(mockUserId)

	// Check if the donations are retrieved successfully
	assert.NoError(t, err)
	assert.NotNil(t, donations)

	mockRepo.AssertExpectations(t)
}

func TestGetAllDonations_Failed(t *testing.T) {
	mockRepo := new(repository.MockDonationRepository)

	// Representing retrieving all donations from the database
	mockUserId := 1
	mockRepo.On("GetAllDonation", mockUserId).Return(nil, assert.AnError)
	donations, err := mockRepo.GetAllDonation(mockUserId)

	// Check if the donations retrieval failed as expected
	assert.Error(t, err)
	assert.Nil(t, donations)

	mockRepo.AssertExpectations(t)
}

func TestCreateDonation_Success(t *testing.T) {
	mockRepo := new(repository.MockDonationRepository)

	// Representing a donation created and retrieved from the database
	mockUserId := 1
	mockDonation := model.Donation{
		ID:         1,
		UserID:     1,
		CampaignID: 1,
		Amount:     50000,
		Message:    "Donation for a cause",
		Status:     "PENDING",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	mockDonationPtr := &mockDonation

	// Representing creating a donation in the database
	mockRepo.On("CreateDonation", mockUserId, mockDonationPtr).Return(mockDonationPtr, nil)
	donationPtr, err := mockRepo.CreateDonation(mockUserId, mockDonationPtr)

	// Check if the donation is created successfully
	assert.NoError(t, err)
	assert.NotNil(t, donationPtr)

	mockRepo.AssertExpectations(t)
}

func TestCreateDonation_Failed(t *testing.T) {
	mockRepo := new(repository.MockDonationRepository)

	// Representing creating a donation in the database
	mockUserId := 1
	mockDonation := model.Donation{
		ID:         1,
		UserID:     1,
		CampaignID: 1,
		Amount:     50000,
		Message:    "Donation for a cause",
		Status:     "PENDING",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	mockDonationPtr := &mockDonation

	mockRepo.On("CreateDonation", mockUserId, mockDonationPtr).Return(nil, assert.AnError)
	donationPtr, err := mockRepo.CreateDonation(mockUserId, mockDonationPtr)

	// Check if the donation creation failed as expected
	assert.Error(t, err)
	assert.Nil(t, donationPtr)

	mockRepo.AssertExpectations(t)
}

func TestUpdateDonation_Success(t *testing.T) {
	mockRepo := new(repository.MockDonationRepository)

	// Representing a donation created and retrieved from the database
	mockUserId := 1
	mockDonation := model.Donation{
		ID:         1,
		UserID:     1,
		CampaignID: 1,
		Amount:     50000,
		Message:    "Donation for a cause",
		Status:     "PENDING",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	mockDonationPtr := &mockDonation

	// Representing updating a donation in the database
	mockRepo.On("UpdateDonation", mockUserId, mockDonationPtr).Return(mockDonationPtr, nil)
	donationPtr, err := mockRepo.UpdateDonation(mockUserId, mockDonationPtr)

	// Check if the donation is updated successfully
	assert.NoError(t, err)
	assert.NotNil(t, donationPtr)

	mockRepo.AssertExpectations(t)
}

func TestUpdateDonation_Failed(t *testing.T) {
	mockRepo := new(repository.MockDonationRepository)

	// Representing a donation created and retrieved from the database
	mockUserId := 1
	mockDonation := model.Donation{
		ID:         1,
		UserID:     1,
		CampaignID: 1,
		Amount:     50000,
		Message:    "Donation for a cause",
		Status:     "PENDING",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	mockDonationPtr := &mockDonation

	// Representing updating a donation in the database
	mockRepo.On("UpdateDonation", mockUserId, mockDonationPtr).Return(nil, assert.AnError)
	donationPtr, err := mockRepo.UpdateDonation(mockUserId, mockDonationPtr)

	// Check if the donation update failed as expected
	assert.Error(t, err)
	assert.Nil(t, donationPtr)

	mockRepo.AssertExpectations(t)
}
