package model

import (
	"errors"
	"time"

)

type Claims struct {
	UserID int     `json:"user_id"`
	Email  string  `json:"email"`
	Exp    float64 `json:"exp"`
}

// Valid method to implement jwt.Claims interface
func (c Claims) Valid() error {
	if time.Now().After(time.Unix(int64(c.Exp), 0)) {
		return errors.New("token has expired")
	}
	return nil
}
