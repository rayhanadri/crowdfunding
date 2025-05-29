package repository

import (
	"errors"
	"time"

	"github.com/rayhanadri/crowdfunding/donation-service/model"
	"gorm.io/gorm"
)

type DonationRepository interface {
	GetAllDonations(user_id int) (*[]model.Donation, error)
	CreateDonation(user_id int, donation *model.Donation) (*model.Donation, error)
	GetDonationByID(user_id int, donationID int) (*model.Donation, error)
	UpdateDonation(user_id int, donation *model.Donation) (*model.Donation, error)
}

type donationRepository struct {
	db *gorm.DB
}

func NewDonationRepository(db *gorm.DB) DonationRepository {
	return &donationRepository{db: db}
}

func (r *donationRepository) GetAllDonations(user_id int) (*[]model.Donation, error) {
	// validate user id
	user := new(model.User)
	if err := r.db.Where("id = ?", user_id).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Fetch all donations for the user
	var donations []model.Donation
	if err := r.db.Where("user_id = ?", user_id).Find(&donations).Error; err != nil {
		return nil, err
	}

	for i := range donations {
		donation := &donations[i]

		// Fetch related campaign
		campaign := new(model.Campaign)
		if err := r.db.Where("id = ?", donation.CampaignID).First(&campaign).Error; err != nil {
			return nil, errors.New("campaign not found for donation")
		}

		// Fetch campaign creator user
		userCreator := new(model.User)
		if err := r.db.Where("id = ?", campaign.UserID).First(&userCreator).Error; err != nil {
			return nil, errors.New("user not found for campaign")
		}

		// Assign nested objects
		campaign.User = *userCreator
		donation.Campaign = *campaign
	}

	return &donations, nil
}

func (r *donationRepository) GetDonationByID(user_id int, donationID int) (*model.Donation, error) {
	// validate user id
	user := new(model.User)
	if err := r.db.Where("id = ?", user_id).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	donation := new(model.Donation)
	if err := r.db.Where("user_id = ? AND id = ?", user_id, donationID).First(donation).Error; err != nil {
		return nil, err
	}

	// get campaign from donation
	campaign := new(model.Campaign)
	if err := r.db.Where("id = ?", donation.CampaignID).First(&campaign).Error; err != nil {
		return nil, errors.New("campaign not found")
	}
	// get campaign creator user
	userCreator := new(model.User)
	if err := r.db.Where("id = ?", campaign.UserID).First(&userCreator).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Assign All Inner Object for json
	campaign.User = *userCreator
	donation.Campaign = *campaign
	donation.User = *user

	return donation, nil
}

func (r *donationRepository) CreateDonation(user_id int, donation *model.Donation) (*model.Donation, error) {
	// validate user id
	user := new(model.User)
	if err := r.db.Where("id = ?", user_id).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// get campaign from donation
	campaign := new(model.Campaign)
	if err := r.db.Where("id = ?", donation.CampaignID).First(&campaign).Error; err != nil {
		return nil, errors.New("campaign not found")
	}

	// chec if campaign ACTIVE
	if campaign.Status != "ACTIVE" {
		return nil, errors.New("campaign not active")
	}

	// get campaign creator user
	userCreator := new(model.User)
	if err := r.db.Where("id = ?", campaign.UserID).First(&userCreator).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// query create db
	if err := r.db.Create(donation).Error; err != nil {
		return nil, err
	}

	// Assign All Inner Object for json
	campaign.User = *userCreator
	donation.Campaign = *campaign
	donation.User = *user

	return donation, nil
}

func (r *donationRepository) UpdateDonation(user_id int, donation *model.Donation) (*model.Donation, error) {
	// validate user id
	user := new(model.User)
	if err := r.db.Where("id = ?", user_id).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// get campaign from donation
	campaign := new(model.Campaign)
	if err := r.db.Where("id = ?", donation.CampaignID).First(&campaign).Error; err != nil {
		return nil, errors.New("campaign not found")
	}
	// get campaign creator user
	userCreator := new(model.User)
	if err := r.db.Where("id = ?", campaign.UserID).First(&userCreator).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// get current date as update date
	donation.UpdatedAt = time.Now()

	// query save donation
	if err := r.db.Save(donation).Error; err != nil {
		return nil, err
	}

	// Assign All Inner Object for json
	campaign.User = *userCreator
	donation.Campaign = *campaign

	return donation, nil
}
