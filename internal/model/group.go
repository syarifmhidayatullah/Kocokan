package model

import "time"

type Group struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	OwnerID        uint      `gorm:"not null" json:"owner_id"`
	Name           string    `gorm:"not null" json:"name"`
	Emoji          string    `json:"emoji"`
	Description    string    `json:"description"`
	NumParticipants int      `gorm:"not null" json:"num_participants"`
	PeriodType     string    `gorm:"not null;default:'monthly'" json:"period_type"` // monthly, weekly, biweekly
	PrizeAmount    int64     `gorm:"not null" json:"prize_amount"`
	TotalRounds    int       `gorm:"not null" json:"total_rounds"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Owner        User          `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Participants []Participant `gorm:"foreignKey:GroupID" json:"participants,omitempty"`
	Rounds       []Round       `gorm:"foreignKey:GroupID" json:"rounds,omitempty"`
}

func (Group) TableName() string { return "koco_groups" }
