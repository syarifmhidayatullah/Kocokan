package model

import "time"

type Participant struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	GroupID   uint      `gorm:"not null" json:"group_id"`
	Name      string    `gorm:"not null" json:"name"`
	Phone     string    `json:"phone"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Participant) TableName() string { return "koco_participants" }
