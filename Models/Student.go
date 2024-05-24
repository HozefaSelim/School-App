package Models

import (
	"time"
)

type Student struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Gender    string    `json:"gender"`
	ClassID   uint      `json:"class_id"`
	Class     Class     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"class"`
	Subjects  []Subject `gorm:"many2many:student_subjects;" json:"subjects"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
