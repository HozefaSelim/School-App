package Models

import (
	"time"
)

type Teacher struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Salary    int       `json:"salary"`
	Gender    string    `json:"gender"`
	Email     string    `json:"email"`
	SchoolID  uint      `json:"school_id"`
	School    School    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:Set NULL;"`
	SubjectID uint      `json:"subject_id"`
	Subject   Subject   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
