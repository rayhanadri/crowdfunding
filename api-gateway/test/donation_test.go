package test

import (
	"testing"
	"time"

	"github.com/rayhanadri/crowdfunding/donation-service/model"
	"github.com/stretchr/testify/assert"

	"github.com/rayhanadri/crowdfunding/api-gateway/repository"
)

func TestGetDonationByID_Success(t *testing.T) {
	mockRepo := new(repository.MockDonationRepository)

	mockDonation := model.Donation{
		ID:          1,
		UserID:      1,
		CampaignID:  1,
		Amount:      50000,
		MessageText: "Donation for a cause",
		Status:      "PENDING",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockDonationPtr := &mockDonation

	// Representing retrieving a donation by ID from the database
	mockRepo.On("GetDonationByID", 1).Return(mockDonationPtr, nil)
	donationPtr, err := mockRepo.GetDonationByID(1)

	// Check if the donation is retrieved successfully
	assert.NoError(t, err)
	assert.NotNil(t, donationPtr)

	mockRepo.AssertExpectations(t)
}

func TestGetDonationByID_Failed(t *testing.T) {
	mockRepo := new(repository.MockDonationRepository)

	// Representing retrieving a donation by ID from the database
	mockRepo.On("GetDonationByID", 1).Return(nil, assert.AnError)
	donationPtr, err := mockRepo.GetDonationByID(1)

	// Check if the donation retrieval failed as expected
	assert.Error(t, err)
	assert.Nil(t, donationPtr)

	mockRepo.AssertExpectations(t)
}

func TestGetAllDonations_Success(t *testing.T) {
	mockRepo := new(repository.MockDonationRepository)

	// Representing donations retrieved from the database
	mockDonations := []model.Donation{
		{
			ID:          1,
			UserID:      1,
			CampaignID:  1,
			Amount:      50000,
			MessageText: "Donation for a cause",
			Status:      "PENDING",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			UserID:      2,
			CampaignID:  2,
			Amount:      50000,
			MessageText: "Donation for a cause",
			Status:      "PENDING",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockDonationsPtr := &mockDonations

	// Representing retrieving all donations from the database
	mockRepo.On("GetAllDonation").Return(mockDonationsPtr, nil)
	donations, err := mockRepo.GetAllDonation()

	// Check if the donations are retrieved successfully
	assert.NoError(t, err)
	assert.NotNil(t, donations)

	mockRepo.AssertExpectations(t)
}

func TestGetAllDonations_Failed(t *testing.T) {
	mockRepo := new(repository.MockDonationRepository)

	// Representing retrieving all donations from the database
	mockRepo.On("GetAllDonation").Return(nil, assert.AnError)
	donations, err := mockRepo.GetAllDonation()

	// Check if the donations retrieval failed as expected
	assert.Error(t, err)
	assert.Nil(t, donations)

	mockRepo.AssertExpectations(t)
}

func TestCreateDonation_Success(t *testing.T) {
	mockRepo := new(repository.MockDonationRepository)

	// Representing a donation created and retrieved from the database
	mockDonation := model.Donation{
		ID:          1,
		UserID:      1,
		CampaignID:  1,
		Amount:      50000,
		MessageText: "Donation for a cause",
		Status:      "PENDING",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockDonationPtr := &mockDonation

	// Representing creating a donation in the database
	mockRepo.On("CreateDonation", mockDonationPtr).Return(mockDonationPtr, nil)
	donationPtr, err := mockRepo.CreateDonation(mockDonationPtr)

	// Check if the donation is created successfully
	assert.NoError(t, err)
	assert.NotNil(t, donationPtr)

	mockRepo.AssertExpectations(t)
}

func TestCreateDonation_Failed(t *testing.T) {
	mockRepo := new(repository.MockDonationRepository)

	// Representing creating a donation in the database
	mockDonation := model.Donation{
		ID:          1,
		UserID:      1,
		CampaignID:  1,
		Amount:      50000,
		MessageText: "Donation for a cause",
		Status:      "PENDING",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockDonationPtr := &mockDonation

	mockRepo.On("CreateDonation", mockDonationPtr).Return(nil, assert.AnError)
	donationPtr, err := mockRepo.CreateDonation(mockDonationPtr)

	// Check if the donation creation failed as expected
	assert.Error(t, err)
	assert.Nil(t, donationPtr)

	mockRepo.AssertExpectations(t)
}

func TestUpdateDonation_Success(t *testing.T) {
	mockRepo := new(repository.MockDonationRepository)

	// Representing a donation created and retrieved from the database
	mockDonation := model.Donation{
		ID:          1,
		UserID:      1,
		CampaignID:  1,
		Amount:      50000,
		MessageText: "Donation for a cause",
		Status:      "PENDING",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockDonationPtr := &mockDonation

	// Representing updating a donation in the database
	mockRepo.On("UpdateDonation", mockDonationPtr).Return(mockDonationPtr, nil)
	donationPtr, err := mockRepo.UpdateDonation(mockDonationPtr)

	// Check if the donation is updated successfully
	assert.NoError(t, err)
	assert.NotNil(t, donationPtr)

	mockRepo.AssertExpectations(t)
}

func TestUpdateDonation_Failed(t *testing.T) {
	mockRepo := new(repository.MockDonationRepository)

	// Representing a donation created and retrieved from the database
	mockDonation := model.Donation{
		ID:          1,
		UserID:      1,
		CampaignID:  1,
		Amount:      50000,
		MessageText: "Donation for a cause",
		Status:      "PENDING",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockDonationPtr := &mockDonation

	// Representing updating a donation in the database
	mockRepo.On("UpdateDonation", mockDonationPtr).Return(nil, assert.AnError)
	donationPtr, err := mockRepo.UpdateDonation(mockDonationPtr)

	// Check if the donation update failed as expected
	assert.Error(t, err)
	assert.Nil(t, donationPtr)

	mockRepo.AssertExpectations(t)
}
