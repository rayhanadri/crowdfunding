package model

import (
	"time"

	user_model "github.com/rayhanadri/crowdfunding/user-service/model" // corrected the import path
)

type Donation struct {
	ID         int             `gorm:"primaryKey" json:"id"`
	UserID     int             `json:"user_id"`
	User       user_model.User `gorm:"foreignKey:UserID" json:"user"` // corrected the import path
	CampaignID int             `json:"campaign_id"`
	// Campaign   Campaign `gorm:"foreignKey:CampaignID" json:"campaign"` // corrected the import path

	Amount    float64   `json:"amount"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Donation) TableName() string {
	return "donations.donations"
}
