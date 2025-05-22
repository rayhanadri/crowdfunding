package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"crowdfund/model"
)

type TransactionRepository interface {
	GetAllTransaction(user_id int) (*[]model.Transaction, error)
	CreateTransaction(user_id int, transaction *model.Transaction) (*model.Transaction, error)
	GetTransactionByID(user_id int, transactionID int) (*model.Transaction, error)
	UpdateTransaction(user_id int, transaction *model.Transaction) (*model.Transaction, error)
	CheckUpdateTransaction(user_id int, transaction *model.Transaction) (*model.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) GetAllTransaction(user_id int) (*[]model.Transaction, error) {
	// validate user id
	user := new(model.User)
	if err := r.db.Where("id = ?", user_id).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	var transactions []model.Transaction
	query := `
		SELECT t.*
		FROM transactions t
		LEFT JOIN donations d
		ON t.donation_id = d.id
		WHERE d.user_id = ?
	`
	if err := r.db.Raw(query, user_id).Scan(&transactions).Error; err != nil {
		return nil, err
	}

	for i := range transactions {
		transaction := &transactions[i]

		// Fetch related donation
		donation := new(model.Donation)
		if err := r.db.Where("id = ?", transaction.DonationID).First(&donation).Error; err != nil {
			return nil, errors.New("donation not found for transaction")
		}

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
		transaction.Donation = *donation
	}

	return &transactions, nil
}

func (r *transactionRepository) CreateTransaction(user_id int, transaction *model.Transaction) (*model.Transaction, error) {
	// validate user id
	user := new(model.User)
	if err := r.db.Where("id = ?", user_id).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// get donation from donation id
	donation := new(model.Donation)
	if err := r.db.Where("id = ?", transaction.DonationID).First(&donation).Error; err != nil {
		return nil, errors.New("donation not found")
	}

	// check if donation user_id match
	if donation.UserID != user_id {
		return nil, errors.New("user not authorized")
	}

	// check if donation with user and donation id already exist
	var existingTransaction model.Transaction
	if err := r.db.Where("user_id = ? AND donation_id = ?", user_id, transaction.DonationID).First(&existingTransaction).Error; err == nil {
		return nil, errors.New("transaction already exists")
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

	//fill transaction description based on campaign title
	transaction.InvoiceDescription = campaign.Title

	//query create db
	if err := r.db.Create(transaction).Error; err != nil {
		return nil, err
	}

	// Assign All Inner Object for json
	campaign.User = *userCreator
	donation.Campaign = *campaign
	donation.User = *user
	transaction.Donation = *donation

	return transaction, nil
}

func (r *transactionRepository) UpdateTransaction(user_id int, transaction *model.Transaction) (*model.Transaction, error) {
	// validate user id
	user := new(model.User)
	if err := r.db.Where("id = ?", user_id).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// get donation from donation id
	donation := new(model.Donation)
	if err := r.db.Where("id = ?", transaction.DonationID).First(&donation).Error; err != nil {
		return nil, errors.New("donation not found")
	}

	// check if donation user_id match
	if donation.UserID != user_id {
		return nil, errors.New("user not authorized")
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
	transaction.UpdatedAt = time.Now()

	// query save transaction
	if err := r.db.Save(transaction).Error; err != nil {
		return nil, err
	}

	// Assign All Inner Object for json
	campaign.User = *userCreator
	donation.Campaign = *campaign
	transaction.Donation = *donation

	return transaction, nil
}

func (r *transactionRepository) GetTransactionByID(user_id int, transactionID int) (*model.Transaction, error) {
	// validate user id
	user := new(model.User)
	if err := r.db.Where("id = ?", user_id).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	transaction := new(model.Transaction)
	if err := r.db.Where("user_id = ? AND id = ?", user_id, transactionID).First(transaction).Error; err != nil {
		return nil, err
	}

	// get donation from donation id
	donation := new(model.Donation)
	if err := r.db.Where("id = ?", transaction.DonationID).First(&donation).Error; err != nil {
		return nil, errors.New("donation not found")
	}

	// check if donation user_id match
	if donation.UserID != user_id {
		return nil, errors.New("user not authorized")
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
	transaction.Donation = *donation

	return transaction, nil
}

func (r *transactionRepository) CheckUpdateTransaction(user_id int, transaction *model.Transaction) (*model.Transaction, error) {
	// validate user id
	user := new(model.User)
	if err := r.db.Where("id = ?", user_id).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// get donation from donation id
	donation := new(model.Donation)
	if err := r.db.Where("id = ?", transaction.DonationID).First(&donation).Error; err != nil {
		return nil, errors.New("donation not found")
	}

	// check if donation user_id match
	if donation.UserID != user_id {
		return nil, errors.New("user not authorized")
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
	transaction.UpdatedAt = time.Now()

	// query save transaction
	if err := r.db.Save(transaction).Error; err != nil {
		return nil, err
	}

	// update donation
	donation.Status = transaction.Status

	// update campaign
	if transaction.Status == "PAID" {
		campaign.CollectedAmount += transaction.Amount

		// save campaign
		if err := r.db.Save(campaign).Error; err != nil {
			return nil, err
		}
	}

	// save donation
	if err := r.db.Save(donation).Error; err != nil {
		return nil, err
	}

	// Assign All Inner Object for json
	campaign.User = *userCreator
	donation.Campaign = *campaign
	transaction.Donation = *donation

	return transaction, nil
}
