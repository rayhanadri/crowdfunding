package model

type Campaign struct {
	ID     int  `gorm:"primaryKey" json:"id"`
	UserID int  `json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`

	Title           string  `json:"title"`
	Description     string  `json:"description"`
	TargetAmount    float64 `json:"target_amount"`
	CollectedAmount float64 `json:"collected_amount"`
	Deadline        string  `json:"deadline"`
	Status          string  `json:"status"`
	Category        string  `json:"category,omitempty"`
	MinDonation     float64 `json:"min_donation,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

func (Campaign) TableName() string {
	return "campaigns"
}
