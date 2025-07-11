package models

import (
	"time"

	"github.com/google/uuid"
)

type Expense struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid" json:"user_id"`
	Amount      float64   `json:"amount"`
	Currency    string    `gorm:"size:3" json:"currency"`
	Category    string    `json:"category"`
	Description string    `gorm:"size:256" json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}
