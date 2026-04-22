package repository

import (
	"github.com/project/kocokan/internal/model"
	"gorm.io/gorm"
)

type RoundRepository struct{ db *gorm.DB }

func NewRoundRepository(db *gorm.DB) *RoundRepository { return &RoundRepository{db} }

func (r *RoundRepository) ListByGroup(groupID uint) ([]model.Round, error) {
	var rs []model.Round
	err := r.db.Where("group_id = ?", groupID).Preload("Winner").
		Order("round_number ASC").Find(&rs).Error
	return rs, err
}

func (r *RoundRepository) FindByID(id uint) (*model.Round, error) {
	var round model.Round
	err := r.db.Preload("Winner").First(&round, id).Error
	return &round, err
}

func (r *RoundRepository) Create(round *model.Round) error {
	return r.db.Create(round).Error
}

func (r *RoundRepository) Save(round *model.Round) error {
	return r.db.Save(round).Error
}

func (r *RoundRepository) WinnerIDs(groupID uint) ([]uint, error) {
	var ids []uint
	err := r.db.Model(&model.Round{}).
		Where("group_id = ? AND winner_id IS NOT NULL", groupID).
		Pluck("winner_id", &ids).Error
	return ids, err
}
