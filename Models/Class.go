package Models

import (
	"time"
)

type Class struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `json:"name"`
	SchoolID  uint      `json:"school_id"`
	School    School    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
