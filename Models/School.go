package Models

import (
	"time"
)

type School struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	City      string    `json:"city"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
