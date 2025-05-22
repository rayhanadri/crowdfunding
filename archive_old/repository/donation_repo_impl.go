package repository

import (
	"github.com/stretchr/testify/mock"

	"crowdfund/model"
)

type MockUserDonationInterface interface {
	GetAllDonation(user_id int) (*[]model.Donation, error)
	CreateDonation(user_id int, donation *model.Donation) (*model.Donation, error)
	GetDonationByID(user_id int, donationID int) (*model.Donation, error)
	UpdateDonation(user_id int, donation *model.Donation) (*model.Donation, error)
}

type MockDonationRepository struct {
	mock.Mock
}

func (m *MockDonationRepository) GetAllDonation(user_id int) (*[]model.Donation, error) {
	args := m.Called(user_id)
	if donations := args.Get(0); donations != nil {
		return donations.(*[]model.Donation), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDonationRepository) CreateDonation(user_id int, donation *model.Donation) (*model.Donation, error) {
	args := m.Called(user_id, donation)
	if donation := args.Get(0); donation != nil {
		return donation.(*model.Donation), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDonationRepository) GetDonationByID(user_id int, donationID int) (*model.Donation, error) {
	args := m.Called(user_id, donationID)
	if donation := args.Get(0); donation != nil {
		return donation.(*model.Donation), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDonationRepository) UpdateDonation(user_id int, donation *model.Donation) (*model.Donation, error) {
	args := m.Called(user_id, donation)
	if donation := args.Get(0); donation != nil {
		return donation.(*model.Donation), args.Error(1)
	}
	return nil, args.Error(1)
}
