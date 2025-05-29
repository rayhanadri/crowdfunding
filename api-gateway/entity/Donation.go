package entity

import "time"

type Donation struct {
	ID         int      `gorm:"primaryKey" json:"id"`
	UserID     int      `json:"user_id"`
	User       User     `gorm:"foreignKey:UserID" json:"user"`
	CampaignID int      `json:"campaign_id"`
	Campaign   Campaign `gorm:"foreignKey:CampaignID" json:"campaign"`

	Amount    float64   `json:"amount"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Donation) TableName() string {
	return "donations"
}

type DonationRequest struct {
	ID         int       `gorm:"primaryKey" json:"id"`
	UserID     int       `json:"user_id"`
	CampaignID int       `json:"campaign_id"`
	Amount     float64   `json:"amount"`
	Message    string    `json:"message"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
