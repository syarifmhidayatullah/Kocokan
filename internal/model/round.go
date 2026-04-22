package model

import "time"

type Round struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	GroupID       uint         `gorm:"not null" json:"group_id"`
	RoundNumber   int          `gorm:"not null" json:"round_number"`
	WinnerID      *uint        `json:"winner_id"`
	DrawnAt       *time.Time   `json:"drawn_at"`
	Notes         string       `json:"notes"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`

	Winner        *Participant `gorm:"foreignKey:WinnerID" json:"winner,omitempty"`
}

func (Round) TableName() string { return "koco_rounds" }
