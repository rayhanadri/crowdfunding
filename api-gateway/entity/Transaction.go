package entity

import "time"

type Transaction struct {
	ID         int      `gorm:"primaryKey" json:"id"`
	DonationID int      `gorm:"not null;index" json:"donation_id"`
	Donation   Donation `gorm:"foreignKey:DonationID" json:"donation"`

	InvoiceID          string    `gorm:"size:255" json:"invoice_id"`
	InvoiceURL         string    `gorm:"size:255" json:"invoice_url"`
	InvoiceDescription string    `gorm:"size:255" json:"invoice_description"`
	PaymentMethod      string    `gorm:"size:50" json:"payment_method"`
	Amount             float64   `gorm:"not null" json:"amount"`
	Status             string    `gorm:"size:50;default:'PENDING'" json:"status"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Transaction) TableName() string {
	return "transactions"
}
